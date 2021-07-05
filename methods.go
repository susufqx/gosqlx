package gosqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/susufqx/gosqlx/util"
)

// Read : find by the options
func Read(ctx context.Context, p PreparerContext, baseModels interface{}, options map[string]interface{}) error {
	return read(ctx, p, baseModels, options, nil, nil, nil, nil)
}

// ReadPageSort : find by order and offset
func ReadPageSort(ctx context.Context, p PreparerContext, baseModels interface{}, options map[string]interface{}, size, offset int, orderKey, orderDire string) error {
	if size < 0 {
		baseModels = []BaseModelInterface{}
		return nil
	}

	return read(ctx, p, baseModels, options, &size, &offset, &orderKey, &orderDire)
}

// Save : if the db model exists, update the content,
// or insert a new to db
func Save(ctx context.Context, p PreparerContext, baseModel BaseModelInterface) error {
	allMap, pkMap, noPkMap := collectKV(ctx, baseModel)
	return save(ctx, p, baseModel.GetTableName(), allMap, noPkMap, pkMap)
}

// Create : create a new, no judge if the model exists
func Create(ctx context.Context, p PreparerContext, baseModel BaseModelInterface) error {
	allMap, _, _ := collectKV(ctx, baseModel)
	return create(ctx, p, baseModel.GetTableName(), allMap)
}

// Update : update the data without judging the model's existance
func Update(ctx context.Context, p PreparerContext, baseModel BaseModelInterface) error {
	_, pkMap, noPkMap := collectKV(ctx, baseModel)
	return update(ctx, p, baseModel.GetTableName(), noPkMap, pkMap)
}

// UpdateMap : map is to record the update key-values
func UpdateMap(ctx context.Context, p PreparerContext, tableName string, qm map[string]interface{}, cm map[string]interface{}) error {
	return update(ctx, p, tableName, cm, qm)
}

// Delete : delete the data by primary keys by default
func Delete(ctx context.Context, p PreparerContext, options ...interface{}) error {
	optionsLen := len(options)
	if optionsLen < 1 {
		return errors.New("parameters error")
	}

	baseModel := options[0]
	bm, ok := baseModel.(BaseModelInterface)
	if !ok {
		return errors.New("parameters error")
	}

	_, pkMap, _ := collectKV(ctx, bm)

	var mapOptions map[string]interface{}
	if optionsLen == 2 {
		mapOptions, ok = options[1].(map[string]interface{})
		if !ok {
			return errors.New("parameters error")
		}

		pkMap = util.MapJoin(pkMap, mapOptions)
	} else if optionsLen > 2 {
		if optionsLen/2 == 0 {
			return errors.New("new pairs of key-value, but got key not value")
		}

		for i := 1; i < optionsLen; i = i + 2 {
			key, ok := options[i].(string)
			if !ok {
				return errors.New("need string, but got others")
			}

			pkMap[key] = options[i+1]
		}
	}

	return delete(ctx, p, bm.GetTableName(), pkMap)
}

func save(ctx context.Context, p PreparerContext, tableName string, allMap, noPkMap, pkMap map[string]interface{}) error {
	allLen := len(allMap)
	noPkLen := len(noPkMap)
	pkLen := len(pkMap)

	keys := make([]string, allLen+pkLen)
	dollarSigns := make([]string, allLen+noPkLen)
	values := make([]interface{}, allLen+noPkLen)
	i := 0
	for key, value := range allMap {
		keys[i] = key
		dollarSigns[i] = "$" + strconv.Itoa(i+1)
		values[i] = value
		i++
	}

	j, k := i, i
	for key := range pkMap {
		keys[j] = key
		j++
	}

	for key, value := range noPkMap {
		dollarSigns[k] = key + "=$" + strconv.Itoa(k+1)
		values[k] = value
		k++
	}

	st := fmt.Sprintf("INSERT INTO %s("+strings.Join(keys[:allLen], ",")+
		") VALUES ("+strings.Join(dollarSigns[:allLen], ",")+
		") ON conflict("+strings.Join(keys[allLen:], ",")+
		") DO UPDATE SET "+strings.Join(dollarSigns[allLen:], ","),
		tableName)

	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, values...)
	return err
}

func create(ctx context.Context, p PreparerContext, tableName string, allMap map[string]interface{}) error {
	keys := make([]string, len(allMap))
	dollarSigns := make([]string, len(allMap))
	values := make([]interface{}, len(allMap))
	i := 0
	for key, value := range allMap {
		keys[i] = key
		dollarSigns[i] = "$" + strconv.Itoa(i+1)
		values[i] = value
		i++
	}

	st := fmt.Sprintf("INSERT INTO %s ("+strings.Join(keys, ",")+") VALUES ("+strings.Join(dollarSigns, ",")+")", tableName)
	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, values...)
	return err
}

func read(ctx context.Context, p PreparerContext, baseModels interface{}, kv map[string]interface{}, size, offset *int, orderStr, orderSort *string) error {
	var tableName = getTableName(ctx, baseModels)
	signs := make([]string, len(kv))
	values := []interface{}{}
	i := 0
	valueCount := 0
	for k, v := range kv {
		if vs, ok := v.([]interface{}); ok {
			tsigns := make([]string, len(vs))
			for k, it := range vs {
				tsigns[k] = "$" + strconv.Itoa(valueCount+1)
				values = append(values, it)
				valueCount++
			}
			signs[i] = k + " IN (" + strings.Join(tsigns, ",") + ")"
		} else {
			signs[i] = k + "=$" + strconv.Itoa(i+1)
			values = append(values, v)
			valueCount++
		}

		i++
	}

	var st string
	if kv == nil {
		st = fmt.Sprintf("SELECT * FROM %s", tableName)
	} else {
		st = fmt.Sprintf("SELECT * FROM %s WHERE "+strings.Join(signs, " AND "), tableName)
	}

	if orderStr != nil && orderSort != nil {
		st = st + " ORDER BY " + *orderStr + " " + *orderSort
	}

	if size != nil {
		st = st + " LIMIT " + strconv.Itoa(*size)
		if offset != nil {
			st = st + " OFFSET " + strconv.Itoa(*offset)
		}
	}

	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return err
	}

	err = stmt.Select(baseModels, values...)
	return err
}

func update(ctx context.Context, p PreparerContext, tableName string, noPkMap, pKmap map[string]interface{}) error {
	noPkLen := len(noPkMap)
	pkLen := len(pKmap)
	signsM := make([]string, noPkLen+pkLen)
	values := []interface{}{}
	i := 0
	for key, value := range noPkMap {
		signsM[i] = key + "=$" + strconv.Itoa(i+1)
		values = append(values, value)
		i++
	}

	var valueCount = i
	for key, value := range pKmap {
		if vs, ok := value.([]interface{}); ok {
			tsigns := make([]string, len(vs))
			for k, it := range vs {
				tsigns[k] = "$" + strconv.Itoa(valueCount+1)
				values = append(values, it)
				valueCount++
			}
			signsM[i] = key + " IN (" + strings.Join(tsigns, ",") + ")"
		} else {
			signsM[i] = key + "=$" + strconv.Itoa(i+1)
			values = append(values, value)
			valueCount++
		}

		i++
	}

	st := fmt.Sprintf("UPDATE %s SET "+strings.Join(signsM[:noPkLen], ",")+" WHERE "+strings.Join(signsM[noPkLen:], " AND "), tableName)
	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, values...)
	return err
}

func delete(ctx context.Context, p PreparerContext, tableName string, kv map[string]interface{}) error {
	signs := make([]string, len(kv))
	values := []interface{}{}
	i := 0
	valueCount := 0
	for k, v := range kv {
		if vs, ok := v.([]interface{}); ok {
			tsigns := make([]string, len(vs))
			for k, it := range vs {
				tsigns[k] = "$" + strconv.Itoa(valueCount+1)
				values = append(values, it)
				valueCount++
			}
			signs[i] = k + " IN (" + strings.Join(tsigns, ",") + ")"
		} else {
			signs[i] = k + "=$" + strconv.Itoa(i+1)
			values = append(values, v)
			valueCount++
		}

		i++
	}

	st := fmt.Sprintf("DELETE FROM %s WHERE "+strings.Join(signs, " AND "), tableName)
	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, values...)
	return err
}

func existByPrimaryKeys(ctx context.Context, p PreparerContext, tableName string, mp map[string]interface{}) (bool, error) {
	signs := make([]string, len(mp))
	values := make([]interface{}, len(mp))
	i := 0
	for key, value := range mp {
		signs[i] = key + "=$" + strconv.Itoa(i+1)
		values[i] = value
		i++
	}

	st := fmt.Sprintf("SELECT * FROM %s WHERE "+strings.Join(signs, " AND "), tableName)
	stmt, err := p.PreparexContext(ctx, st)
	if err != nil {
		return false, err
	}

	var baseModel BaseModelInterface
	row := stmt.QueryRowxContext(ctx, values...)
	err = row.StructScan(baseModel)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func collectKV(ctx context.Context, baseModel BaseModelInterface) (allMap, pkMap, noPkMap map[string]interface{}) {
	typ := reflect.TypeOf(baseModel).Elem()  // get type definition
	val := reflect.ValueOf(baseModel).Elem() // get val elements
	num := typ.NumField()

	allMap = map[string]interface{}{}
	pkMap = map[string]interface{}{}
	noPkMap = map[string]interface{}{}

	for i := 0; i < num; i++ {
		var k string

		fd := typ.Field(i)
		vl := val.Field(i)
		tagDB := fd.Tag.Get("db")
		tagOthers := fd.Tag.Get("others")
		pKey := false

		if tagDB == "" {
			k = util.CamelCaseToSnackCase(fd.Name)
			allMap[k] = vl.Interface()
		} else {
			k = tagDB
			allMap[k] = vl.Interface()
		}

		if tagOthers != "" {
			vs := strings.Split(tagOthers, ";")
			for _, vss := range vs {
				if vss == "pKey" {
					pKey = true
				}
			}

		}

		if pKey {
			pkMap[k] = allMap[k]
		} else {
			noPkMap[k] = allMap[k]
		}
	}

	return
}

func getTableName(ctx context.Context, baseModels interface{}) string {
	getType := reflect.TypeOf(baseModels).Elem()
	if getType.Kind() == reflect.Slice {
		getType = getType.Elem()
	}

	vp := reflect.New(getType)

	met, ok := getType.MethodByName("GetTableName")
	if !ok {
		return util.CamelCaseToSnackCase(getType.String())
	}

	r := met.Func.Call([]reflect.Value{reflect.Indirect(vp)})
	return r[0].String()
}

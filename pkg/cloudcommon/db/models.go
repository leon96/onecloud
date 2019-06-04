// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"fmt"

	"yunion.io/x/log"
	"yunion.io/x/pkg/utils"
	"yunion.io/x/sqlchemy"
)

var globalTables map[string]IModelManager

func RegisterModelManager(modelMan IModelManager) {
	if globalTables == nil {
		globalTables = make(map[string]IModelManager)
	}
	mustCheckModelManager(modelMan)
	log.Infof("Register model %s", modelMan.Keyword())
	globalTables[modelMan.Keyword()] = modelMan
}

func mustCheckModelManager(modelMan IModelManager) {
	allowedTags := map[string][]string{
		"create": {"required", "optional", "domain_required", "domain_optional", "admin_required", "admin_optional"},
		"search": {"user", "domain", "admin"},
		"get":    {"user", "domain", "admin"},
		"list":   {"user", "domain", "admin"},
		"update": {"user", "domain", "admin"},
	}
	for _, col := range modelMan.TableSpec().Columns() {
		tags := col.Tags()
		for tagName, allowedValues := range allowedTags {
			v, ok := tags[tagName]
			if !ok {
				continue
			}
			if !utils.IsInStringArray(v, allowedValues) {
				msg := fmt.Sprintf("model manager %s: column %s has invalid tag %s:\"%s\", expecting %v",
					modelMan.KeywordPlural(), col.Name(), tagName, v, allowedValues)
				panic(msg)
			}
		}
	}
}

func CheckSync(autoSync bool) bool {
	log.Infof("Start check database ...")
	examinedTables := make(map[string]bool)
	allDropFKSqls := make([]string, 0)
	allSqls := make([]string, 0)
	for modelName, modelMan := range globalTables {
		log.Infof("# check table of model %s", modelName)
		tableSpec := modelMan.TableSpec()
		if _, ok := examinedTables[tableSpec.Name()]; ok {
			continue
		}
		examinedTables[tableSpec.Name()] = true
		dropFKSqls := tableSpec.DropForeignKeySQL()
		if len(dropFKSqls) > 0 {
			allDropFKSqls = append(allDropFKSqls, dropFKSqls...)
		}
		sqls := tableSpec.SyncSQL()
		if len(sqls) > 0 {
			allSqls = append(allSqls, sqls...)
		}
	}
	allSqls = append(allDropFKSqls, allSqls...)
	if len(allSqls) > 0 {
		if autoSync {
			err := commitSqlDIffs(allSqls)
			if err == nil {
				return true
			} else {
				log.Errorln(err)
			}
		}
		for _, sql := range allSqls {
			fmt.Println(sql)
		}
		log.Fatalf("Database not in sync!")
		return false
	} else {
		log.Infof("Database is in SYNC!!!")
		return true
	}
}

func GetModelManager(keyword string) IModelManager {
	modelMan, ok := globalTables[keyword]
	if ok {
		return modelMan
	} else {
		return nil
	}
}

func commitSqlDIffs(sqls []string) error {
	db := sqlchemy.GetDB()

	for _, sql := range sqls {
		log.Infof("Exec %s", sql)
		_, err := db.Exec(sql)
		if err != nil {
			log.Errorf("Exec sql failed %s\n%s", sql, err)
			return err
		}
	}
	return nil
}

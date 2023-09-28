package data_manager

import "jsondb/internal/files_manager"

// type CollectionManager struct {
// }

func CreateCollection(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	path := MAIN_DIR_PATH + "/" + projectName
	return files_manager.CreateFile(path, kwargs["name"].(string))
}

func RenamCollection(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	path := MAIN_DIR_PATH + "/" + projectName + "/" + kwargs["names"].([]interface{})[0].(string) + ".json"
	return files_manager.RenameFile(path, kwargs["names"].([]interface{})[1].(string))
}

func DeleteCollection(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	path := MAIN_DIR_PATH + "/" + projectName + "/" + kwargs["name"].(string) + ".json"
	return files_manager.DeleteFile(path)
}

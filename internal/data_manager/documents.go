package data_manager

import "jsondb/internal/files_manager"

// type DocumentsManager struct {
// }

func CreateDocument(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	path := MAIN_DIR_PATH + "/" + projectName + "/" + kwargs["collection_name"].(string) + ".json"

	fetchedData, err := files_manager.ReadFile(path)

	if err != nil {
		return err
	}

	fetchedData = append(fetchedData, kwargs["doc"].(map[string]interface{}))

	err = files_manager.WriteFile(path, fetchedData)

	return err
}

func UpdateDocument(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {

	path := MAIN_DIR_PATH + "/" + projectName + "/" + kwargs["collection_name"].(string) + ".json"

	_, err := files_manager.ReadFile(path)

	if err != nil {
		return err
	}
	return nil
}

func DeleteDocument(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	return nil
}

func GetAllDocuments(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	return nil
}

func GetOneDocument(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	return nil
}

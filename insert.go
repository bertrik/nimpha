package main

import (
	"bytes"
    "fmt"
	"log"
	// "errors"	
	//"io"
	"encoding/json"
	"strconv"
	"net/http"
	"rsocket_json_requests"
)

func insert(w http.ResponseWriter, r *http.Request) {
	

	key_id := r.URL.Query().Get("key_id")
	ikey_id, err := strconv.Atoi(key_id)
	if err != nil {
		log.Fatal(err)
	}
	var result mem_row
	var coll []mem_row
	result.Key_id = ikey_id
	result.Table_name = r.URL.Query().Get("table")
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

    var p interface {}
	err = dec.Decode(&p)
	result.Document = p
	result.Parsed_Document = p.(map[string]interface{})

	//Check if IndexRow is full. Then create another and append.Otherwise, just append to the mem_table and ++ the counter.
	//The next One should be rotational list of available servers
	//create keep alive
	//aftter that, create method to update INDEX TABLES through the servers and create WRITE AHEAD LOG to be shared among the servers and order the indexes according to the request.

	coll = append(coll, result)
	index_it := get_wal(&coll)
	index_row := it.Index_rows[index_it] 
	json_data, err := json.Marshal(coll)

	fmt.Println("JSON_DATA")
	fmt.Println(coll)


    if err != nil {
        log.Fatal(err)
	}
	
	response, err := http.Post("http://" +  index_row.Instance_ip + ":" + index_row.Instance_port + "/" + index_row.Instance_name +  "/insert_worker?table=" + result.Table_name + "&key_id=" + strconv.Itoa(result.Key_id), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
	 	log.Fatal(err)
	}

	fmt.Println(response)

	fmt.Fprintf(w,"Success")
}

func insert_worker(w http.ResponseWriter, r *http.Request) {
	
	key_id := r.URL.Query().Get("key_id")
	ikey_id, err := strconv.Atoi(key_id)
	if err != nil {
		log.Fatal(err)
	}
	var result mem_row
	result.Key_id = ikey_id
	//result.Document = r.URL.Query().Get("document")
	result.Table_name = r.URL.Query().Get("table")
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	

    var p []mem_row
	err = dec.Decode(&p)
	
	fmt.Println("Output-----------")
	fmt.Println(p)

	result.Document = p[0].Document
	result.Parsed_Document = p[0].Document.(map[string]interface{})
	mt.Rows = append(mt.Rows, result)

	fmt.Println(mt)
	//Check if IndexRow is full. Then create another and append.Otherwise, just append to the mem_table and ++ the counter.
	//The next One should be rotational list of available servers
	//create keep alive



	fmt.Fprintf(w, "Success")
}


func insert_rsocket(payload interface{})  interface{} {
	

	payload_content := make(map[string]interface{})
	myString := payload.(string)
	json.Unmarshal([]byte(myString), &payload_content)
	key_id := payload_content["key_id"].(string)
	ikey_id, err := strconv.Atoi(key_id)


	if err != nil {
		log.Fatal(err)
	}
	var result mem_row
	var coll []mem_row
	result.Key_id = ikey_id
	result.Table_name = payload_content["table"].(string)
	intermediate_inteface := payload_content["body"].(map[string]interface{})
	
	fmt.Println(intermediate_inteface)
	fmt.Println(payload_content)

	result.Document = intermediate_inteface
	result.Parsed_Document = intermediate_inteface


	//Check if IndexRow is full. Then create another and append.Otherwise, just append to the mem_table and ++ the counter.
	//The next One should be rotational list of available servers
	//create keep alive
	//aftter that, create method to update INDEX TABLES through the servers and create WRITE AHEAD LOG to be shared among the servers and order the indexes according to the request.

	coll = append(coll, result)
	index_it := get_wal_rsocket(&coll)
	index_row := it.Index_rows[index_it] 

	fmt.Println(coll)
	var param interface{}
	param = map[string]interface{}{
		"key_id": strconv.Itoa(result.Key_id),
		"table":  result.Table_name,
		"body": coll,
	}
	fmt.Println(param)
    if err != nil {
        log.Fatal(err)
	}
	
	jsonParam, _ := json.Marshal(param)
	_port, _ := strconv.Atoi(index_row.Instance_port)
	rsocket_json_requests.RequestConfigs(index_row.Instance_ip, _port)
	response, err := rsocket_json_requests.RequestJSON("/" + index_row.Instance_name +  "/insert_worker", string(jsonParam))
	if (err != nil){
		fmt.Println(err)
	}

	fmt.Println(response)

	//fmt.Fprintf(w,"Success")
	return response
}

func insert_worker_rsocket(payload interface{}) interface{} {

	payload_content := make(map[string]interface{})
	myString := payload.(string)
	json.Unmarshal([]byte(myString), &payload_content)

	key_id := payload_content["key_id"].(string)
	ikey_id, err := strconv.Atoi(key_id)
	if err != nil {
		log.Fatal(err)
	}
	var result mem_row
	result.Key_id = ikey_id
	//result.Document = r.URL.Query().Get("document")
	result.Table_name = payload_content["table"].(string)
	
	fmt.Println(payload_content["body"])
	fmt.Println("-----------INSERT_WORKER_1-----------")
	fmt.Println(payload_content["body"].([]interface{}))


	intermediate_inteface := payload_content["body"].([]interface{})

	json_rows_bytes, _ := json.Marshal(intermediate_inteface)
	
	fmt.Println(intermediate_inteface)
	reader := bytes.NewReader(json_rows_bytes)

	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	

    var p []mem_row
	err = dec.Decode(&p)
	
	fmt.Println("-----------Output-----------")
	fmt.Println(p)

	result.Document = p[0].Document
	result.Parsed_Document = p[0].Document.(map[string]interface{})
	mt.Rows = append(mt.Rows, result)

	//fmt.Println(mt)

	return "Success"
}
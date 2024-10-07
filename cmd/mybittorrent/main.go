package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
)

var _ = json.Marshal

func decode(str string, index int) (interface{}, int, error) {
    if index >= len(str) {
        return nil, index, fmt.Errorf("Invalid string index %d", index)
    }

    switch {
    case str[index] == 'l':
        return decodeList(str, index)
    case str[index] == 'i':
        return decodeNumber(str, index)
    case str[index] >= '0' && str[index] <= '9':
        return decodeString(str, index)
    case str[index] == 'd':
        return decodeDict(str, index)
    default:
        return nil, index, fmt.Errorf("Unexpected error value %q", str[index])
    }
}

func decodeString(str string, index int) (interface{}, int, error) {
    var strLenStr string
    for i := index; i < len(str); i++ {
        if str[i] == ':' {
            index = i + 1
            break
        }
        strLenStr += (string(str[i]))
    }
    strLen, err := strconv.Atoi(strLenStr)
    if err != nil {
        fmt.Println("Error while converting string to number")
        return nil, index, err
    }
    retStr := str[index : index+strLen]
    index += strLen
    return retStr, index, nil
}

func decodeNumber(str string, index int) (interface{}, int, error) {
    var numStr string
    for i := index + 1; i < len(str); i++ {
        index = i
        if str[i] == 'e' {
            break
        }
        numStr += (string(str[i]))
    }
    decodedNumber, err := strconv.Atoi(numStr)
    if err != nil {
        fmt.Println("Error while converting string to number")
        return nil, index, err
    }
    return decodedNumber, index + 1, nil
}

func decodeList(str string, index int) (interface{}, int, error) {
    index++

    list := make([]interface{}, 0)

    for {
        if index >= len(str) {
            return nil, index, fmt.Errorf("bad list")
        }

        if str[index] == 'e' {
            break
        }

        var x interface{}

        x, i, err := decode(str, index)
        if err != nil {
            return nil, index, err
        }
        index = i
        list = append(list, x)
    }

    return list, index + 1, nil
}

func decodeDict(str string, index int) (interface{}, int, error) {
    index++

    dict := make(map[string]interface{})
    for {
        if index >= len(str) {
            return nil, index, fmt.Errorf("bad dict")
        }

        if str[index] == 'e' {
            break
        }

        key, i, err := decode(str, index)
        index = i
        if err != nil {
            return nil, i, err
        }
        if index >= len(str) {
            return nil, index, fmt.Errorf("bad dict")
        }
        value, i, err := decode(str, index)
        dict[key.(string)] = value
        index = i
    }
    return dict, index + 1, nil
}

func main() {
    command := os.Args[1]

    switch command {
    case "decode":
        bencodedValue := os.Args[2]

        decodedString, _, err := decode(bencodedValue, 0)
        if err != nil {
            fmt.Println("Error while decoding the string ", err)
            return
        }
        jsonOutput, _ := json.Marshal(decodedString)

        fmt.Println(string(jsonOutput))
    case "info":
        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Println("Error reading the torrent file ", err)
            os.Exit(1)
        }

        d, _, err := decodeDict(string(data),  0)
        if err != nil {
            fmt.Println("Error decoding the torrent file ", err)
            os.Exit(1)
        }
        fmt.Printf("Tracker URL: %v\n", d.(map[string]interface{})["announce"])
        info, ok := d.(map[string]interface{})["info"]
        if info == nil || !ok {
            fmt.Println("No info section in the torrent file")
            return
        }

        fmt.Printf("Length: %v\n", info.(map[string]interface{})["length"])
    default:
        fmt.Println("Unknown command: " + command)
        os.Exit(1)
    }
}

package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "unicode"
)

var _ = json.Marshal
var err error

func handleBencodedList(bencodedString string) ([]interface{}, error) {
    var decodedList []interface{};

    ind := 1
    for ind < len(bencodedString)-1 {
        var currString string
        if bencodedString[ind] == 'i' {
            for bencodedString[ind] != 'e' {
                currString += string(bencodedString[ind]);
                ind++
            }
            currString += "e"
            ind++
        } else {
            var currLen int
            var nextStringLen int
            for i:= ind; i<len(bencodedString); i++ {
                if bencodedString[i] == ':' {
                    currLen = i-ind
                    nextStringLen, err = strconv.Atoi(bencodedString[ind:ind+currLen])
                }
            }
            currString = bencodedString[ind:ind+currLen+1]
            ind += currLen+1
            currString +=  bencodedString[ind:ind+nextStringLen];
            ind+= nextStringLen
        }
        decodedString ,err := handleBencode(currString)
        if err != nil {
            fmt.Println("Error while decoding the list")
            return nil, err
        }
        decodedList = append(decodedList, decodedString);
    }
    return decodedList, nil
}

func handleBencode(bencodedString string) (interface{}, error) {
    var firstColenIndex =-1

    if bencodedString[0] == 'l' {
        decodedList, err :=  handleBencodedList(bencodedString)
        return decodedList, err
    }

    for i:=0; i<len(bencodedString); i++ {
        if bencodedString[i] == ':' {
            firstColenIndex = i
            break
        }
    }

    // decode a number
    if firstColenIndex == -1 && bencodedString[0]=='i' {
        numString := bencodedString[1:len(bencodedString)-1]
        decodedNumber, err := strconv.Atoi(numString)
        if err != nil {
            return "", fmt.Errorf("Input string does not represent a number")
        }
        return decodedNumber , nil
    }

    stringLen := bencodedString[:firstColenIndex]

    if !unicode.IsDigit(rune(stringLen[0])) {
        return "", fmt.Errorf("bencoded string length can only be a number")
    }

    encodedString :=bencodedString[firstColenIndex+1:]

    return encodedString, nil
}

func main() {

    command := os.Args[1]

    if command == "decode" {
        bencodedValue := os.Args[2]

        decodedString, err := handleBencode(bencodedValue)
        if err != nil {
            fmt.Println("Error while decoding the string ", err)
            return
        }
        jsonOutput, _ := json.Marshal(decodedString)

        if string(jsonOutput) == "null" {
            fmt.Println("[]")
            return
        }

        fmt.Println(string(jsonOutput))

    } else {
        fmt.Println("Unknown command: " + command)
        os.Exit(1)
    }
}

package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "unicode"
)

var _ = json.Marshal

func handleBencode(bencodedString string) (interface{}, error) {
    var firstColenIndex int =-1 

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

        fmt.Println(string(jsonOutput))

    } else {
        fmt.Println("Unknown command: " + command)
        os.Exit(1)
    }
}

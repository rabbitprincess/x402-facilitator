package trc20

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/rabbitprincess/x402-facilitator/utils"
	"golang.org/x/crypto/sha3"
)

type ABI struct {
	Constructor *Method
	Methods     map[string]*Method
	Events      map[string]*Event
}

func (a *ABI) Pack(name string, params ...interface{}) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("constructors are not supported yet")
	}
	if a.Methods[name] == nil {
		return nil, fmt.Errorf("method not found")
	}
	m := a.Methods[name]
	var value []byte
	value = append(value, m.SigId()...)
	value = append(value, m.Inputs.Pack(params...)...)
	return value, nil
}

func (a *ABI) PackParams(name string, params ...interface{}) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("constructors are not supported yet")
	}
	if a.Methods[name] == nil {
		return nil, fmt.Errorf("method not found")
	}
	m := a.Methods[name]
	var value []byte
	value = append(value, m.Inputs.Pack(params...)...)
	return value, nil
}

type Method struct {
	Name    string
	Const   bool
	Inputs  Arguments
	Outputs Arguments
}

func (m *Method) SigId() []byte {
	//function foo(uint32 a, int b)    =    "foo(uint32,int256)"
	types := make([]string, len(m.Inputs))
	for i, v := range m.Inputs {
		types[i] = v.Type
	}
	functionStr := fmt.Sprintf("%v(%v)", m.Name, strings.Join(types, ","))
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write([]byte(functionStr))
	return keccak256.Sum(nil)[:4]
}

type Event struct {
	Name      string
	Anonymous bool
	Inputs    Arguments
}

type Arguments []Argument

type Argument struct {
	Name    string
	Type    string
	Indexed bool // indexed is only used by events
}

func (arg Arguments) Pack(params ...any) []byte {
	if len(arg) != len(params) {
		return nil
	}

	var value []byte
	var i = 0
	for _, v := range arg {
		p := params[i]
		i++
		switch v.Type {
		case "uint256", "uint128", "uint64", "uint32", "uint", "int256", "int128", "int64", "int32", "int":
			va := reflect.ValueOf(p).Interface().(*big.Int)
			value = append(value, PaddedBigBytes(U256(va), 32)...)
		case "string":
			//packBytesSlice([]byte(reflectValue.String()), reflectValue.Len())
			va := new(big.Int).SetBytes([]byte(reflect.ValueOf(p).String())) //
			value = append(value, PaddedBigBytes(U256(va), 32)...)
		case "address":
			addr := reflect.ValueOf(p).String()
			//addrByte, _ := hex.DecodeString(addr)

			va := new(big.Int).SetBytes(utils.RemoveZeroHex(addr))
			value = append(value, PaddedBigBytes(U256(va), 32)...)
		}
	}
	return value
}

const TRC20ABI = `[
    {
        "constant": false,
        "inputs": [
            {
                "name": "spender",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "approve",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function",
        "signature": "0x095ea7b3"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0x18160ddd"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "from",
                "type": "address"
            },
            {
                "name": "to",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "transferFrom",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function",
        "signature": "0x23b872dd"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "spender",
                "type": "address"
            },
            {
                "name": "addedValue",
                "type": "uint256"
            }
        ],
        "name": "increaseAllowance",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function",
        "signature": "0x39509351"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "owner",
                "type": "address"
            }
        ],
        "name": "balanceOf",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0x70a08231"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "spender",
                "type": "address"
            },
            {
                "name": "subtractedValue",
                "type": "uint256"
            }
        ],
        "name": "decreaseAllowance",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function",
        "signature": "0xa457c2d7"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "to",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function",
        "signature": "0xa9059cbb"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "owner",
                "type": "address"
            },
            {
                "name": "spender",
                "type": "address"
            }
        ],
        "name": "allowance",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0xdd62ed3e"
    },
    {
        "inputs": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "symbol",
                "type": "string"
            },
            {
                "name": "decimals",
                "type": "uint8"
            },
            {
                "name": "cap",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "constructor",
        "signature": "constructor"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "from",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "to",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Transfer",
        "type": "event",
        "signature": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "owner",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "spender",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Approval",
        "type": "event",
        "signature": "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0x06fdde03"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0x95d89b41"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [
            {
                "name": "",
                "type": "uint8"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function",
        "signature": "0x313ce567"
    }
]`

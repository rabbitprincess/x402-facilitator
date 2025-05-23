package types

type Signer func(digest []byte) (signature []byte, err error)

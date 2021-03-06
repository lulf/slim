/*
 * Copyright 2019, Ulf Lilleengen
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */
package api

type Message struct {
	Offset  int64
	Payload []byte
}

func NewMessage(offset int64, payload []byte) *Message {
	return &Message{
		Offset:  offset,
		Payload: payload,
	}
}

// Copyright (c) 2012 - Cloud Instruments Co., Ltd.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"bytes"
	"io"
	"os"
)

// LoggerFromConfigAsFile creates logger with config from file. File should contain valid seelog xml.
func LoggerFromConfigAsFile(fileName string) (LoggerInterface, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf, err := configFromReader(file)
	if err != nil {
		return nil, err
	}

	return createLoggerFromConfig(conf)
}

// LoggerFromConfigAsBytes creates a logger with config from bytes stream. Bytes should contain valid seelog xml.
func LoggerFromConfigAsBytes(data []byte) (LoggerInterface, error) {
	conf, err := configFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return createLoggerFromConfig(conf)
}

// LoggerFromConfigAsString creates a logger with config from a string. String should contain valid seelog xml.
func LoggerFromConfigAsString(data string) (LoggerInterface, error) {
	return LoggerFromConfigAsBytes([]byte(data))
}

// LoggerFromWriterWithMinLevel creates a proxy logger that uses io.Writer as the
// receiver with minimal level = minLevel.
//
// All messages with level more or equal to minLevel will be written to output and
// formatted using the default seelog format.
//
// Can be called for usage with non-Seelog systems
func LoggerFromWriterWithMinLevel(output io.Writer, minLevel LogLevel) (LoggerInterface, error) {
	constraints, err := newMinMaxConstraints(minLevel, CriticalLvl)
	if err != nil {
		return nil, err
	}

	dispatcher, err := newSplitDispatcher(defaultformatter, []interface{}{output})
	if err != nil {
		return nil, err
	}

	conf, err := newConfig(constraints, make([]*logLevelException, 0), dispatcher, syncloggerTypeFromString, nil)
	if err != nil {
		return nil, err
	}

	return createLoggerFromConfig(conf)
}

// LoggerFromCustomReceiver creates a proxy logger that uses a CustomReceiver as the
// receiver.
//
// All messages will be sent to the specified custom receiver without additional
// formatting ('%Msg' format is used).
//
// Check CustomReceiver, RegisterReceiver for additional info.
//
// NOTE 1: CustomReceiver.AfterParse is only called when a receiver is instantiated
// by the config parser while parsing config. So, if you are not planning to use the
// same CustomReceiver for both proxying (via LoggerFromCustomReceiver call) and
// loading from config, just leave AfterParse implementation empty.
//
// NOTE 2: Unlike RegisterReceiver, LoggerFromCustomReceiver takes an already initialized
// instance that implements CustomReceiver. So, fill it with data and perform any initialization
// logic before calling this func and it won't be lost.
//
// So:
// * RegisterReceiver takes value just to get the reflect.Type from it and then
// instantiate it as many times as config is reloaded.
//
// * LoggerFromCustomReceiver takes value and uses it without modification and
// reinstantiation, directy passing it to the dispatcher tree.
func LoggerFromCustomReceiver(receiver CustomReceiver) (LoggerInterface, error) {
	constraints, err := newMinMaxConstraints(TraceLvl, CriticalLvl)
	if err != nil {
		return nil, err
	}

	output, err := newCustomReceiverDispatcherByValue(msgonlyformatter, receiver)
	if err != nil {
		return nil, err
	}
	dispatcher, err := newSplitDispatcher(msgonlyformatter, []interface{}{output})
	if err != nil {
		return nil, err
	}

	conf, err := newConfig(constraints, make([]*logLevelException, 0), dispatcher, syncloggerTypeFromString, nil)
	if err != nil {
		return nil, err
	}

	return createLoggerFromConfig(conf)
}

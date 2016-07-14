// Copyright (c) 2016, German Neuroinformatics Node (G-Node)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package util

import (
	"bytes"
	"net/smtp"
	"strings"
	"text/template"
)

// EmailDispatcher defines an interface for e-mail dispatch.
type EmailDispatcher interface {
	Send(recipient []string, message []byte) error
}

// EmailConfig contains all information required for e-mail dispatch via smtp.
type EmailConfig struct {
	Identity   string
	Dispatcher string
	Password   string
	Host       string
	Port       string
}

type emailDispatcher struct {
	conf EmailConfig
	send func(string, smtp.Auth, string, []string, []byte) error
}

// Send sets up authentication for e-mail dispatch via smtp and invokes the objects send function.
func (e *emailDispatcher) Send(recipient []string, content []byte) error {
	addr := e.conf.Host + ":" + e.conf.Port
	auth := smtp.PlainAuth(e.conf.Identity, e.conf.Dispatcher, e.conf.Password, e.conf.Host)
	return e.send(addr, auth, e.conf.Dispatcher, recipient, content)
}

// NewSmtpSendMailDispatcher returns an instance of emailDispatcher
// using smtp.SendMail as send function.
func NewSmtpSendMailDispatcher(conf EmailConfig) EmailDispatcher {
	return &emailDispatcher{conf, smtp.SendMail}
}

// MakePlainEmailTemplate returns a bytes.Buffer containing a standard e-mail
func MakePlainEmailTemplate(from string, to []string, subj string, messageBody string) *bytes.Buffer {
	const emailTemplate = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

{{ .Body }}
`
	var doc bytes.Buffer

	content := &struct {
		From    string
		To      string
		Subject string
		Body    string
	}{
		from,
		strings.Join(to, ", "),
		subj,
		messageBody,
	}
	t := template.New("emailTemplate")
	t, err := t.Parse(emailTemplate)
	if err != nil {
		panic("Error parsing e-mail template")
	}
	err = t.Execute(&doc, content)
	if err != nil {
		panic("Error executing e-mail template")
	}
	return &doc
}

// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"
)

// Warning holds a short message that's meant to alert about system events.
// There'll only ever be one Warning with the same message, and it can be
// silenced for a while before repeating. After a (supposedly longer) while
// it'll go away on its own (unless it recurs).
type Warning struct {
	Message     string        `json:"message"`
	FirstAdded  time.Time     `json:"first-added"`
	LastAdded   time.Time     `json:"last-added"`
	LastShown   time.Time     `json:"last-shown,omitempty"`
	ExpireAfter time.Duration `json:"expire-after,omitempty"`
	RepeatAfter time.Duration `json:"repeat-after,omitempty"`
}

type jsonWarning struct {
	Warning
	ExpireAfter string `json:"expire-after,omitempty"`
	RepeatAfter string `json:"repeat-after,omitempty"`
}

type WarningsOptions struct {
	// All means return all warnings, instead of only the un-okayed ones.
	All bool
}

// Warnings returns the list of un-okayed warnings.
func (client *Client) Warnings(opts WarningsOptions) ([]*Warning, error) {
	var jws []*jsonWarning
	q := make(url.Values)
	if opts.All {
		q.Add("select", "all")
	}
	_, err := client.DoSync(&RequestInfo{
		Method: "GET",
		Path:   "/v1/warnings",
		Query:  q,
	}, &jws)

	ws := make([]*Warning, len(jws))
	for i, jw := range jws {
		ws[i] = &jw.Warning
		ws[i].ExpireAfter, _ = time.ParseDuration(jw.ExpireAfter)
		ws[i].RepeatAfter, _ = time.ParseDuration(jw.RepeatAfter)
	}

	return ws, err
}

type warningsAction struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// Okay asks the server to silence the warnings that would have been returned by
// Warnings at the given time.
func (client *Client) Okay(t time.Time) error {
	var body bytes.Buffer
	var op = warningsAction{Action: "okay", Timestamp: t}
	if err := json.NewEncoder(&body).Encode(op); err != nil {
		return err
	}
	_, err := client.DoSync(&RequestInfo{
		Method: "POST",
		Path:   "/v1/warnings",
		Body:   &body,
	}, nil)
	return err
}

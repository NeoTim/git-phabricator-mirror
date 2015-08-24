/*
Copyright 2015 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package request defines the internal representation of a review request.
package request

import (
	"encoding/json"
	"github.com/google/git-phabricator-mirror/mirror/repository"
)

const notesRef = "refs/notes/devtools/reviews"

// Ref defines the git-notes ref that we expect to contain review requests.
var Ref = repository.NotesRef(notesRef)

// Request is
type Request struct {
	ReviewRef   string   `json:"reviewRef,omitempty"`
	TargetRef   string   `json:"targetRef"`
	Requester   string   `json:"requester,omitempty"`
	Reviewers   []string `json:"reviewers,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Parse parses a review request from a git note.
func Parse(note repository.Note) (Request, error) {
	bytes := []byte(note)
	var request Request
	err := json.Unmarshal(bytes, &request)
	// TODO(ojarjur): If "requester" is not set, then use git-blame to fill it in.
	return request, err
}

// ParseAllValid takes collection of git notes and tries to parse a review
// request from each one. Any notes that are not valid review requests get
// ignored, as we expect the git notes to be a heterogenous list, with only
// some of them being review requests.
func ParseAllValid(notes []repository.Note) []Request {
	var requests []Request
	for _, note := range notes {
		request, err := Parse(note)
		if err == nil && request.TargetRef != "" {
			requests = append(requests, request)
		}
	}
	return requests
}

// Write writes a review request as a JSON-formatted git note.
func (request *Request) Write() (repository.Note, error) {
	bytes, err := json.Marshal(request)
	return repository.Note(bytes), err
}

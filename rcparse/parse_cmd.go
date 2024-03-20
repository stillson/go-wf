/*
 * Copyright (c) 2024. Christopher Stillson <stillson@gmail.com>
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package rcparse

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	Start = iota
	InDoubleq
	InSingleq
	InWord
	InSpace
)

func ParseCmd(cmd string) (string, []string, error) {
	outslice := []string{}
	lastSlice := []rune{}
	cmd = strings.Trim(cmd, " \t\n")
	rcmd := []rune(cmd)
	state := Start

	for i, r := range rcmd {

		switch state {
		case Start:
			{
				switch {
				case unicode.IsSpace(r):
					state = InSpace
				case r == '"':
					state = InDoubleq
				case r == '\'':
					state = InSingleq
				case unicode.IsPrint(r):
					{
						// this would include a space char, but we have covered that
						state = InWord
						lastSlice = append(lastSlice, r)
					}
				default:
					return "", nil, fmt.Errorf("parsing Error at %d in %v", i, cmd)
				}

			}
		case InDoubleq:
			{
				switch {
				case r == '"':
					{
						outslice = append(outslice, string(lastSlice))
						// maybe copy this?
						lastSlice = []rune{}
						state = Start
					}
				case unicode.IsPrint(r):
					lastSlice = append(lastSlice, r)
				default:
					return "", nil, fmt.Errorf("parsing Error at %d in %v", i, cmd)
				}
			}
		case InSingleq:
			{
				switch {
				case r == '\'':
					{
						outslice = append(outslice, string(lastSlice))
						// maybe copy this?
						lastSlice = []rune{}
						state = Start
					}
				case unicode.IsPrint(r):
					lastSlice = append(lastSlice, r)
				default:
					return "", nil, fmt.Errorf("parsing Error at %d in %v", i, cmd)
				}

			}
		case InWord:
			{
				switch {
				case unicode.IsSpace(r):
					{
						outslice = append(outslice, string(lastSlice))
						// maybe copy this?
						lastSlice = []rune{}
						state = InSpace

					}
				case unicode.IsPrint(r):
					// don't embed a single or double quote in a word
					lastSlice = append(lastSlice, r)
				default:
					return "", nil, fmt.Errorf("parsing Error at %d in %v", i, cmd)
				}
			}
		case InSpace:
			{
				switch {
				case unicode.IsSpace(r):
					state = InSpace
				case r == '"':
					state = InDoubleq
				case r == '\'':
					state = InSingleq
				case unicode.IsPrint(r):
					{
						lastSlice = append(lastSlice, r)
						state = InWord
					}
				default:
					return "", nil, fmt.Errorf("parsing Error at %d in %v", i, cmd)

				}

			}
		}

	}

	if len(lastSlice) != 0 {
		outslice = append(outslice, string(lastSlice))
	}

	return outslice[0], outslice[1:], nil
}

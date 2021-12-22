/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

type Config struct {
	EscapeQueryString bool   `json:"escapeQueryString"`
	DateFormatTZ      string `json:"dateFormatTZ"`
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) HasTimeZone() bool {
	return len(cfg.DateFormatTZ) > 0
}

func (cfg *Config) GetTimeZone() string {
	//t := time.Now()
	//zone, offset := t.Zone()
	//fmt.Println(zone, offset)
	return cfg.DateFormatTZ
}

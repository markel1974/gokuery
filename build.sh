TARGET=./src/version/build.go
CURRVERSION=$(sed -n 's/const BuildVersion = "\(.*\)/\1/p' < $TARGET) CURRVERSION="${CURRVERSION%\"}"
CURRVERSION=$((CURRVERSION + 1))

echo "Current version: ${CURRVERSION}"

cat > $TARGET <<EOF
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

package version

// Code auto generated
const BuildVersion = "${CURRVERSION}"
const BuildDate = "$(date +"%Y-%m-%dT%H:%M:%S%z")"
EOF

APPNAME=$(basename "${PWD}")
echo "Building ${APPNAME}"
go build -o "${APPNAME}"

APPVERSION=$(./"${APPNAME}" -v 2>&1 | grep "${APPNAME}" | cut -f2 -d " ")
echo "Current version $APPVERSION"

ZIPNAME="${APPNAME}.${APPVERSION}.zip"
echo "Creating zipfile ${ZIPNAME}"
zip "${ZIPNAME}" "${APPNAME}"

echo "Done"

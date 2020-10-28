#!/usr/bin/env bash
. ./scripts/version.sh

updateVersion () {
   local v=$1  
   shift
   local yamls=("$@")
   for i in "${yamls[@]}";
   do
       filename=`basename $i`
       case $filename in
       Chart.yaml|swagger.yaml|main.go)
          # replace version
          sed -i "" -E -e "s/\s*version\s*: v.*/version: $v/" $i
          ;;
       *)
          # replace tag
          sed -i "" -E -e "s/\s*tag\s*:.*/tag: $v/" $i
          ;;
       esac
   done
}

update_files=(
   "./main/main.go"
   "./swagger-ui/swagger.yaml"
)
updateVersion "$sheet_crud_version" "${update_files[@]}"

red=$'\e[1;31m'
grn=$'\e[1;32m'
white=$'\e[0m'

echo "versions updated to:"
echo "$grn""sheet-CRUD version:" "$red""$sheet_crud_version""$white"
echo "$grn""files updated:" "$red""$update_files""$white"
#!/bin/bash
program=$1
#some statements
echo "Detecting max input size for $program program by Riley Miranda"
return_code=127
varx=x 
auto_input=x
input_length=1
#loop
while [ 1 -eq 1 ]
do
	$program $auto_input;
	return_code=$?;
	if [ $return_code -eq 127 ]; then
		echo $program does not exist
	elif [[ $return_code -eq 0 ]]; then
		echo $auto_input
		echo Executed successfully with input length = $input_length
		auto_input="$auto_input$varx"
		((input_length=input_length+1))
	else
		((input_length=input_length-1))
		echo $auto_input
		echo Program failed
		echo Max input length = $input_length
		echo Riley Miranda
		break;
	fi
done
#end loop
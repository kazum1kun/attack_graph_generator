s/^([0-9]*),"([^"]*)","([^"]*)",.*$/\t$1 \[label="$1",shape=$3\];/
s/^([0-9]*),"([^"]*)","([^"]*)"$/\t$1 \[label="$1",shape=$3\];/
s/OR/diamond/
s/AND/ellipse/
s/LEAF/box/
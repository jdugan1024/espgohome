BEGIN { 
    print "package espgohome";
    print "import \"fmt\"";
    print "import \"google.golang.org/protobuf/proto\"";
    print "//go:generate stringer -type=MessageID";
    print "type MessageID uint64";
    print "const (";
} 
/message/ { msg = $2 }
/option \(id\)/ { 
    id = $4
    sub(/;/, "", id); 
    printf("\t%sID MessageID = %d\n", msg, id)
    messages[msg] = id
}
END {
    print ")"

    printf("func decodeMessage(raw []byte, msgType MessageID) (proto.Message, error) {\n");
    printf("\tswitch msgType {\n");
    for (m in messages) {
	    printf("\tcase %sID:\n", m);
		printf("\t\tresp := &%s{}\n", m);
		printf("\t\terr := proto.Unmarshal(raw, resp)\n", m);
		printf("\t\treturn resp, err\n");
    }
	printf("\tdefault:\n");
    printf("\t\terr := fmt.Errorf(\"unsupported message: %%d\", msgType)\n");
	printf("\t\treturn nil, err\n");
    printf("\t}\n");
    printf("}\n")
}
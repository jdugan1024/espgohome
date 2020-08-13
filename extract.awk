BEGIN { 
    print "package espgohome";
    print "";
    print "//go:generate stringer -type=MessageID";
    print "type MessageID uint64";
    print "const (";
} 
/message/ { msg = $2 }
/option \(id\)/ { 
    id = $4
    sub(/;/, "", id); 
    printf("\t%sID MessageID = %d\n", msg, id)
}
END { print ")" }
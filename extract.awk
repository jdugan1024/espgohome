BEGIN { print "package espgohome\n\nconst (" } 
/message/ { msg = $2 }
/option \(id\)/ { 
    id = $4
    sub(/;/, "", id); 
    printf("\t%sID = %d\n", msg, id) 
}
END { print ")" }
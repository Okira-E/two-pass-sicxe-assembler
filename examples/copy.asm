. This is a program to copy a file from one location to another.
. The file is copied from the location specified by the START
. statement to the location specified by the END statement.
. The file is copied one record at a time.
. The file is copied until an end of file is encountered.
. The file is copied to the location specified by the RETADR
. Buffer size is 4096 bytes.
. This is a comment because it starts with a period.
COPY    START   1000
FIRST   STL 	RETADR
        LDB 	#LENGTH
        BASE    LENGTH
CLOOP   +JSUB   RDREC
        LDA 	LENGTH
        COMP    #0
        JEQ 	ENDFIL
        +JSUB	WRREC
        J   	CLOOP
ENDFIL  LDA 	EOF
        STA 	BUFFER
        LDA 	#3
        STA 	LENGTH
        +JSUB	WRREC
        J   	@RETADR
EOF     BYTE    C`EOF`
RETADR  RESW    1
LENGTH  RESW    1
BUFFER  RESB    4096
RDREC   CLEAR   X
        CLEAR   A
        CLEAR   S
        +LDT	#4096
RLOOP   TD  	INPUT
        JEQ 	RLOOP
        RD  	INPUT
        COMPR   A, S
        JEQ 	EXIT
        STCH    BUFFER, X
        TIXR    T
        JLT 	RLOOP
EXIT    STX 	LENGTH
        RSUB
INPUT   BYTE    X`F1`
WRREC   CLEAR   X
        LDT 	LENGTH
WLOOP   TD  	OUTPUT
        JEQ 	WLOOP
        LDCH    BUFFER, X
        WD  	OUTPUT
        TIXR    T
        JLT 	WLOOP
        RSUB
OUTPUT  BYTE    X`05`
        END 	FIRST

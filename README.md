# mini PDF editor

## Story

Don't want to upload personal PDF files to online pdf editor website ?
Don't have Acrobat Pro ?

Clone it, run it locally, and edit any pdf you want in this selft hosted pdf editor.

### Features

Very basic, once the program runs, you can :
- upload a pdf file
- edit it 
- - add text, only in arial for now
- - draw (handwriting, or to sign pages)
- download it back.

### Usage : 

```bash
# get the code
git clone ... miniPDFeditor
cd miniPDFeditor

# compile it
go mod tidy
go build

# run it (werserver)
./miniPDFeditor

## Then, in a browser, go to :
http://localhost:8080

# Now use it !
```
Gorage
===

Write File:
---

 #create hash 
 #check if file with hash exists?
  ##yes?
    get file uuid from content table
  ##no?
    store file into filesystem with hash as file name and first x parts of hash as folder name
    create uuid and make entry in content table with
        - uuid
        - create date
        - delete date
        - file size
        - mime type
        - hash
#file entry
    make entry into filetable with:
        - uuid
        - content uuid
        - original-name
        - create date
        - delete date
#return file object

Read File
--- 

#check if file with uuid exists and is not deleted
    ##no 
        return error file not found
#get entry from file table
#enrich with data from content table
#enrich with content from fs
#return file object

Delete File
---

#check if file with uuid exists
    ##no 
        return error file not found
#mark file entry as deleted
#get content uuid from file entry
#check if other files with same content uuid exists
    ##yes
        do nothing
    ##no
        delete file from fs
        mark content table entry as deleted

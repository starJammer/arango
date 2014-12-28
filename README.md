arango
======

## Intent

I'm not a go developer and wanted to use arango for a project on my own.
I decided to write my own driver for using Arango via its REST API.
I did take a look at https://github.com/diegogub/aranGO but I decided
that I wanted to write my own to do more go programming.

## Installing

Run `go get github.com/starJammer/arango`

## Features

* Connect to the database arango server and get info about the current database
* Switch databases
* Create, drop edge and document collections
* Save, Update, Replace, Delete documents
* Save, Update, Replace, Delete edges
* Retrieve document via id only, NO searching by example or AQL queries yet.
* Retrieve documents via simple by example queries

## Upcoming Features


## Usage

I'll add more usage notes about different ways to do things.

    //THE BELOW CODE IS JUST AN EXAMPLE AND WON'T WORK IF YOU
    //TRY AND COMPILE IT

    import ar "github.com/starJammer/arango"
    import "fmt"

    func main() {
        
        //will connect to the _system database by default
        db, err := ar.Conn( "http://localhost:8529" )

        if err != nil {
            //Any errors returned are of the type ArangoError
            aerr, ok := err.(ar.ArangoError)
            if !ok {
                fmt.Println( "Should not get here since if there is an error, its type should be ar.ArangoError" )
            }
            //See errors.go for more info on the ArangoError type
            fmt.Printf(
                "Something went wrong bro....check arango and make sure it's up.\n" + 
                "HTTP CODE : %d\n" + 
                "Arango Err Number : %d\n" + 
                "IsError : %v\n" + 
                "ErrorMessage : %s\n",
                aerr.Code,
                aerr.ErrorNum,
                aerr.IsError,
                aerr.ErrorMessage,
            )
            return
        }

        //You can also use the following two versions
        db, err := ar.ConnDb( "http://localhost:8529", "database_name" )

        //Use this in case you enabled users/passwords
        db, err := ar.ConnDbUserPassword( "http://localhost:8529", "database_name", "username", "password" )
    
        //You can also include the user/password in the host address like this
        db, err := ar.Conn( "http://username:password@localhost:8529" )

        //You can also use HTTPS
        ar.AllowBadSslCerts = true //set this to true if your development certs are not "official looking"
        db, err := ar.Conn( "https://username:password@localhost:8529" )


        //switch to a databes you want to use if you didn't specify the database name when connecting
        db.UseDatabase( "another_database" )

        c, err := db.CreateDocumentCollection( "things" )

        //OR you can do the following for an edge collection
        c, err := db.CreateEdgeCollection( "things" )

        if err != nil {
            fmt.Println( "Things broke again bro, you're failing. I swear it's not my code!")
        }

        type TestDocument struct {
            ar.DocumentImplementation //embed this if you want to get easy access to the Id, Key, and Rev of the documents you save

            MyField string  //Will be save in the database as "MyField" : "value"
            MyField2 bool `json:"my_field"` //Will be saved as "my_field" : "value"
        }

        var testDoc TestDocument
        var testDoc2 = new(TestDocument)

        testDoc.MyField = "some string"
        testDoc.MyField2 = true //just cuz

        //Must pass in a pointer to the document
        err = c.Save( &testDoc )
        err = c.Save( testDoc2 )

        if err != nil {
            //it should be fine though
        }

        //Info about the document from arango
        fmt.Println( testDoc.Id() )
        fmt.Println( testDoc.Key() )
        fmt.Println( testDoc.Rev() )

        testDoc.MyField = "Updating to some other value"

        //notice we pass in the id first and then a pointer to the document
        c.Update( testDoc.Id(), &testDoc )
        
        //the revision will change since we updated it
        fmt.Println( testDoc.Rev() )

    }

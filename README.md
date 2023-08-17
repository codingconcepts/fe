# fe
Export functions out of databases and into code

### Workflow

``` sh
$ fe [database] [language] args
```

So, for example, to create Go code from the functions in a Postgres database:
``` sh
$ fe postgres go --url 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable'
```

### Select statement

``` json
{
   "version":150001,
   "stmts":[
      {
         "stmt":{
            "SelectStmt":{
               "targetList":[
                  {
                     "ResTarget":{
                        "val":{
                           "FuncCall":{
                              "funcname":[
                                 {
                                    "String":{
                                       "sval":"count"
                                    }
                                 }
                              ],
                              "agg_star":true,
                              "funcformat":"COERCE_EXPLICIT_CALL",
                              "location":10
                           }
                        },
                        "location":10
                     }
                  }
               ],
               "fromClause":[
                  {
                     "RangeVar":{
                        "relname":"person",
                        "inh":true,
                        "relpersistence":"p",
                        "location":26
                     }
                  }
               ],
               "whereClause":{
                  "A_Expr":{
                     "kind":"AEXPR_OP",
                     "name":[
                        {
                           "String":{
                              "sval":"="
                           }
                        }
                     ],
                     "lexpr":{
                        "ColumnRef":{
                           "fields":[
                              {
                                 "String":{
                                    "sval":"date_of_birth"
                                 }
                              }
                           ],
                           "location":41
                        }
                     },
                     "rexpr":{
                        "ColumnRef":{
                           "fields":[
                              {
                                 "String":{
                                    "sval":"d"
                                 }
                              }
                           ],
                           "location":57
                        }
                     },
                     "location":55
                  }
               },
               "limitOption":"LIMIT_OPTION_DEFAULT",
               "op":"SETOP_NONE"
            }
         },
         "stmt_len":58
      }
   ]
}
```

### Insert statement

``` json
{
   "version":150001,
   "stmts":[
      {
         "stmt":{
            "InsertStmt":{
               "relation":{
                  "relname":"person",
                  "inh":true,
                  "relpersistence":"p",
                  "location":15
               },
               "cols":[
                  {
                     "ResTarget":{
                        "name":"full_name",
                        "location":23
                     }
                  },
                  {
                     "ResTarget":{
                        "name":"date_of_birth",
                        "location":34
                     }
                  },
                  {
                     "ResTarget":{
                        "name":"country",
                        "location":49
                     }
                  }
               ],
               "selectStmt":{
                  "SelectStmt":{
                     "valuesLists":[
                        {
                           "List":{
                              "items":[
                                 {
                                    "ColumnRef":{
                                       "fields":[
                                          {
                                             "String":{
                                                "sval":"full_name"
                                             }
                                          }
                                       ],
                                       "location":70
                                    }
                                 },
                                 {
                                    "ColumnRef":{
                                       "fields":[
                                          {
                                             "String":{
                                                "sval":"date_of_birth"
                                             }
                                          }
                                       ],
                                       "location":81
                                    }
                                 },
                                 {
                                    "ColumnRef":{
                                       "fields":[
                                          {
                                             "String":{
                                                "sval":"country"
                                             }
                                          }
                                       ],
                                       "location":96
                                    }
                                 }
                              ]
                           }
                        }
                     ],
                     "limitOption":"LIMIT_OPTION_DEFAULT",
                     "op":"SETOP_NONE"
                  }
               },
               "override":"OVERRIDING_NOT_SET"
            }
         },
         "stmt_len":104
      }
   ]
}
```
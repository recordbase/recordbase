# recordbase

Recordbase API and Client

Distributed simple record storage based on Raft protocol. 
Can handle million of records with the strong consistency.
Encrypts all data on storage and transfer.
Could run on different availability zones and provides eviction policy of data.
Recommended to install on machines with SSDs.

# Client

Recordbase client is the HA client supporting connection to multiple instances in the cluster.
The initial connection point is used to pick up the first available data node, request all nodes in the cluster and
to connect to all nodes. Therefore, initial connection point could have one node, but client would be connected to all of them.
Client has embedded load balancer based on gRPC protocol.
API endpoints are protected by JWT tokens.

Discovery of new joined nodes on runtime is not available in the current version, if in case your cluster increased by certain nodes, you
still need to restart clients (applications) to get benefits of newcomers. But it also could be the situation when newcomers are still not
synchronized by Raft protocol, therefore please restart clients when all updates are propagated.

```
rb := recordbase.NewClient(context.Background(), "127.0.0.1:5555,127.0.0.1:7777", MY_API_KEY, tlsConfigOpt)
defer rb.Destroy()
```

# Data Model

* Recordbase is a multi tenant database. The purpose of having tenants is to accommodate situations of M&A and keep all entities in one database.
* Recordbase has a flat keyspace, where PrimaryKey is the unique identifier for each record. Keyspace could be spitted on different zones by forming each separate cluster. 
PrimaryKeys always represent Base62 numbers with some starting point and range. There is allocation mechanism to differentiate key ranges between clusters. 
In case if recordbase is using to store PII data with certain GDPR policy, then management of key ranges would require routing and allocation on upper level.
* Each record in Recordbase contains attributes, tags, columns, maps and files.
   
The description of record objects:
* Attribute - foreign key of the record forming together by the attribute name and attribute value, not necessary unique
* Tag - foreign key of the record, not necessary unique
* Column - entry of key and value keeping inside of the record
* Map - entry of key and value keeping outside of the record
* File - entry of key (file name) keeping inside of record and value (file payload) outside of the record

Names could participate in formatting internal keys, therefore they can have limitations on characters.
Names limitations (a-zA-Z0-9PSC):
* Attribute name 
* Tag
* Map key

PSC (printable special characters):
```
ch == '-' || ch == '_' ||
ch == '`' || ch == '~' || ch == '|' ||
ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '&' || ch == '*' || ch == '=' || ch == '+' ||
ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == '[' || ch == ']' ||
ch == ';' || ch == '<' || ch == '>' || ch == '?' 
```




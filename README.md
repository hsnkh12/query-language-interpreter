# query-language-interpreter

## Descriptioon 
A combination of lexical and syntax anaylzer built from scratch to interpret a custom database query language, built using Go programing language

## Sample query language
create project 'project_name'; 

delete project 'project_name';

rename project 'project_name' 'new_name';


create collection 'collection_name';

delete collection 'collection_name';

rename collection 'collection_name' 'new_name';



add into 'collection_name' doc('attr1': 'value', 'attr2': 'value', 'attr3': doc( 'attr4' : 'value'));

get from 'collection_name' attrs('attr1', 'attr2') where('attr1' == 'attr2' || 'attr2' > 'attr4');

get one from 'collection_name' attrs() where();

update from 'collection_name' set('attr1': 'new_value', 'attr2': 'new_value') where('attr1' == 'attr2' || 'attr2' > 'attr4' && 'attr2' > 'attr4');

delete from 'collection_name' where();



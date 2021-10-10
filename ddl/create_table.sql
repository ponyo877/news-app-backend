CREATE TABLE articleTBL(
   id          serial    NOT NULL,
   title       text,
   URL         text,
   image       text,
   updateDate  timestamp,
   click       integer,
   siteID      integer,
   PRIMARY KEY (id)
);



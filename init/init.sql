CREATE TABLE public.books (
                              id    int NOT NULL,
                              title character varying(255) NOT NULL,
                              author character varying(255) NOT NULL,
                              price numeric(5,2) NOT NULL
);
ALTER TABLE public.books ADD CONSTRAINT  uni_key PRIMARY KEY (id);

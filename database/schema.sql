--
-- PostgreSQL database dump
--

-- Dumped from database version 17.0 (Debian 17.0-1.pgdg120+1)
-- Dumped by pg_dump version 17.1 (Ubuntu 17.1-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: payment_method; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.payment_method AS ENUM (
    'bill',
    'pix',
    'credit_card'
);


ALTER TYPE public.payment_method OWNER TO postgres;

--
-- Name: transaction_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_status AS ENUM (
    'pending',
    'cancelled',
    'expired',
    'approved'
);


ALTER TYPE public.transaction_status OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: branchdailystats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.branchdailystats (
    id integer NOT NULL,
    tenant_id character varying NOT NULL,
    branch_id character varying NOT NULL,
    date date NOT NULL,
    total_boletos integer DEFAULT 0,
    total_pagos integer DEFAULT 0,
    valor_emitido double precision DEFAULT 0.00,
    valor_recebido double precision DEFAULT 0.00,
    boletos_cancelados integer DEFAULT 0,
    valor_cancelado double precision DEFAULT 0.00,
    boletos_atrasados integer DEFAULT 0,
    total_dias_atraso double precision DEFAULT 0.00
);


ALTER TABLE public.branchdailystats OWNER TO postgres;

--
-- Name: branchdailystats_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.branchdailystats_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.branchdailystats_id_seq OWNER TO postgres;

--
-- Name: branchdailystats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.branchdailystats_id_seq OWNED BY public.branchdailystats.id;


--
-- Name: customermonthlystats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customermonthlystats (
    id integer NOT NULL,
    customer_document_number character varying NOT NULL,
    tenant_id character varying NOT NULL,
    month date NOT NULL,
    total_boletos integer DEFAULT 0,
    total_pagos integer DEFAULT 0,
    valor_emitido double precision DEFAULT 0.00,
    valor_recebido double precision DEFAULT 0.00,
    boletos_atrasados integer DEFAULT 0,
    total_dias_atraso double precision DEFAULT 0.00
);


ALTER TABLE public.customermonthlystats OWNER TO postgres;

--
-- Name: customermonthlystats_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.customermonthlystats_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.customermonthlystats_id_seq OWNER TO postgres;

--
-- Name: customermonthlystats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.customermonthlystats_id_seq OWNED BY public.customermonthlystats.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: transaction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    bank_slip_uuid uuid,
    status public.transaction_status,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    due_date timestamp without time zone NOT NULL,
    total numeric NOT NULL,
    customer_document_number character varying NOT NULL,
    tenant_id character varying NOT NULL,
    branch_id character varying NOT NULL,
    payment_method public.payment_method,
    CONSTRAINT check_payment_method CHECK ((payment_method = ANY (ARRAY['bill'::public.payment_method, 'pix'::public.payment_method, 'credit_card'::public.payment_method]))),
    CONSTRAINT check_status CHECK ((status = ANY (ARRAY['pending'::public.transaction_status, 'cancelled'::public.transaction_status, 'expired'::public.transaction_status, 'approved'::public.transaction_status])))
);


ALTER TABLE public.transaction OWNER TO postgres;

--
-- Name: branchdailystats id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branchdailystats ALTER COLUMN id SET DEFAULT nextval('public.branchdailystats_id_seq'::regclass);


--
-- Name: customermonthlystats id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customermonthlystats ALTER COLUMN id SET DEFAULT nextval('public.customermonthlystats_id_seq'::regclass);


--
-- Name: branchdailystats branchdailystats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branchdailystats
    ADD CONSTRAINT branchdailystats_pkey PRIMARY KEY (id);


--
-- Name: customermonthlystats customermonthlystats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customermonthlystats
    ADD CONSTRAINT customermonthlystats_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transaction transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_pkey PRIMARY KEY (id);


--
-- Name: customermonthlystats unique_customer_month; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customermonthlystats
    ADD CONSTRAINT unique_customer_month UNIQUE (customer_document_number, month);


--
-- Name: branchdailystats unique_tenant_branch_date; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branchdailystats
    ADD CONSTRAINT unique_tenant_branch_date UNIQUE (tenant_id, branch_id, date);


--
-- PostgreSQL database dump complete
--


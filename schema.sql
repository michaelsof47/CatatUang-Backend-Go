--
-- PostgreSQL database dump
--

\restrict 2erUcp1Mks6sxG2gVxL9WZNQ0U4o4QsnaRWLuIgzxb5g2OLMoI7tN7GTH0NRnEr

-- Dumped from database version 14.21 (Postgres.app)
-- Dumped by pg_dump version 14.21 (Postgres.app)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: enum_AccountType_account_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public."enum_AccountType_account_status" AS ENUM (
    'active',
    'inactive',
    'suspended'
);


--
-- Name: enum_AccountTypes_account_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public."enum_AccountTypes_account_status" AS ENUM (
    'active',
    'inactive',
    'suspended'
);


--
-- Name: enum_User_reward_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public."enum_User_reward_status" AS ENUM (
    'bronze',
    'silver',
    'gold',
    'platinum'
);


--
-- Name: enum_Users_reward_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public."enum_Users_reward_status" AS ENUM (
    'bronze',
    'silver',
    'gold',
    'platinum'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: AccountType; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."AccountType" (
    id integer NOT NULL,
    account_status public."enum_AccountType_account_status" DEFAULT 'active'::public."enum_AccountType_account_status" NOT NULL
);


--
-- Name: AccountType_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."AccountType_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: AccountType_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."AccountType_id_seq" OWNED BY public."AccountType".id;


--
-- Name: Balance; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Balance" (
    id bigint NOT NULL,
    balance_amount bigint DEFAULT 0 NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: Balance_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public."Balance" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Balance_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: Category; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Category" (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255),
    url_image text,
    user_id uuid NOT NULL
);


--
-- Name: Category_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Category_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Category_id_seq" OWNED BY public."Category".id;


--
-- Name: DetailPlannerBook; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."DetailPlannerBook" (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    total_amount double precision NOT NULL,
    planner_book_id integer NOT NULL
);


--
-- Name: PercentagePlannerBook; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."PercentagePlannerBook" (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    percentage_amount double precision,
    total_amount double precision NOT NULL,
    planner_book_id integer NOT NULL
);


--
-- Name: PercentagePlannerBook_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."PercentagePlannerBook_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: PercentagePlannerBook_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."PercentagePlannerBook_id_seq" OWNED BY public."DetailPlannerBook".id;


--
-- Name: PercentagePlannerBook_id_seq1; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."PercentagePlannerBook_id_seq1"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: PercentagePlannerBook_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."PercentagePlannerBook_id_seq1" OWNED BY public."PercentagePlannerBook".id;


--
-- Name: PlannerBook; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."PlannerBook" (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    target_start timestamp with time zone NOT NULL,
    target_end timestamp with time zone NOT NULL,
    target_amount integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    user_id uuid NOT NULL
);


--
-- Name: PlannerBook_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."PlannerBook_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: PlannerBook_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."PlannerBook_id_seq" OWNED BY public."PlannerBook".id;


--
-- Name: Transaction; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Transaction" (
    id bigint NOT NULL,
    name text NOT NULL,
    transaction_date timestamp with time zone NOT NULL,
    amount integer NOT NULL,
    merchant_name text,
    item_price numeric NOT NULL,
    total_price numeric NOT NULL,
    category_id bigint,
    balances_id bigint NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: Transaction_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public."Transaction" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Transaction_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: User; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."User" (
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    reward_status public."enum_User_reward_status" DEFAULT 'bronze'::public."enum_User_reward_status" NOT NULL,
    url_user_image text,
    email character varying(255) NOT NULL,
    phone character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    account_type_id integer DEFAULT 1 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL
);


--
-- Name: AccountType id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."AccountType" ALTER COLUMN id SET DEFAULT nextval('public."AccountType_id_seq"'::regclass);


--
-- Name: Category id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Category" ALTER COLUMN id SET DEFAULT nextval('public."Category_id_seq"'::regclass);


--
-- Name: DetailPlannerBook id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."DetailPlannerBook" ALTER COLUMN id SET DEFAULT nextval('public."PercentagePlannerBook_id_seq"'::regclass);


--
-- Name: PercentagePlannerBook id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PercentagePlannerBook" ALTER COLUMN id SET DEFAULT nextval('public."PercentagePlannerBook_id_seq1"'::regclass);


--
-- Name: PlannerBook id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PlannerBook" ALTER COLUMN id SET DEFAULT nextval('public."PlannerBook_id_seq"'::regclass);


--
-- Name: AccountType AccountType_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."AccountType"
    ADD CONSTRAINT "AccountType_pkey" PRIMARY KEY (id);


--
-- Name: Balance Balance_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Balance"
    ADD CONSTRAINT "Balance_pkey" PRIMARY KEY (id);


--
-- Name: Category Category_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Category"
    ADD CONSTRAINT "Category_pkey" PRIMARY KEY (id);


--
-- Name: DetailPlannerBook PercentagePlannerBook_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."DetailPlannerBook"
    ADD CONSTRAINT "PercentagePlannerBook_pkey" PRIMARY KEY (id);


--
-- Name: PercentagePlannerBook PercentagePlannerBook_pkey1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PercentagePlannerBook"
    ADD CONSTRAINT "PercentagePlannerBook_pkey1" PRIMARY KEY (id);


--
-- Name: PlannerBook PlannerBook_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PlannerBook"
    ADD CONSTRAINT "PlannerBook_pkey" PRIMARY KEY (id);


--
-- Name: Transaction Transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Transaction"
    ADD CONSTRAINT "Transaction_pkey" PRIMARY KEY (id);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: Balance balance_user_id_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Balance"
    ADD CONSTRAINT balance_user_id_unique UNIQUE (user_id);


--
-- Name: User user_email_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT user_email_unique UNIQUE (email);


--
-- Name: Balance Balance_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Balance"
    ADD CONSTRAINT "Balance_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON UPDATE CASCADE;


--
-- Name: Category Category_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Category"
    ADD CONSTRAINT "Category_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: DetailPlannerBook PercentagePlannerBook_planner_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."DetailPlannerBook"
    ADD CONSTRAINT "PercentagePlannerBook_planner_book_id_fkey" FOREIGN KEY (planner_book_id) REFERENCES public."PlannerBook"(id);


--
-- Name: PercentagePlannerBook PercentagePlannerBook_planner_book_id_fkey1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PercentagePlannerBook"
    ADD CONSTRAINT "PercentagePlannerBook_planner_book_id_fkey1" FOREIGN KEY (planner_book_id) REFERENCES public."PlannerBook"(id);


--
-- Name: PlannerBook PlannerBook_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."PlannerBook"
    ADD CONSTRAINT "PlannerBook_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON UPDATE CASCADE;


--
-- Name: Transaction Transaction_balances_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Transaction"
    ADD CONSTRAINT "Transaction_balances_id_fkey" FOREIGN KEY (balances_id) REFERENCES public."Balance"(id);


--
-- Name: Transaction Transaction_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Transaction"
    ADD CONSTRAINT "Transaction_category_id_fkey" FOREIGN KEY (category_id) REFERENCES public."Category"(id);


--
-- Name: Transaction Transaction_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Transaction"
    ADD CONSTRAINT "Transaction_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: User User_account_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_account_type_id_fkey" FOREIGN KEY (account_type_id) REFERENCES public."AccountType"(id) ON UPDATE CASCADE;

--
-- Name: AccountType; schema: public; Owner: -
--

INSERT INTO public."AccountType" (id, account_status)
VALUES (1, 'active')
ON CONFLICT (id) DO NOTHING;

SELECT setval('public."AccountType_id_seq"', (SELECT max(id) FROM public."AccountType"));

--
-- PostgreSQL database dump complete
--

\unrestrict 2erUcp1Mks6sxG2gVxL9WZNQ0U4o4QsnaRWLuIgzxb5g2OLMoI7tN7GTH0NRnEr


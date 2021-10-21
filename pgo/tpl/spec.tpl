--
-- PostgreSQL database dump
--

-- Dumped from database version 13.4
-- Dumped by pg_dump version 14.0

-- Started on 2021-10-21 14:06:22 EDT

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
-- TOC entry 10 (class 2615 OID 16385)
-- Name: spec; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA spec;


ALTER SCHEMA spec OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 211 (class 1259 OID 16475)
-- Name: applications; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.applications (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.applications OWNER TO postgres;

--
-- TOC entry 210 (class 1259 OID 16461)
-- Name: channels; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.channels (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.channels OWNER TO postgres;

--
-- TOC entry 212 (class 1259 OID 16489)
-- Name: clusterdeployments; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.clusterdeployments (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.clusterdeployments OWNER TO postgres;

--
-- TOC entry 208 (class 1259 OID 16433)
-- Name: configs; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.configs (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.configs OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 16720)
-- Name: klusterletaddonconfigs; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.klusterletaddonconfigs (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.klusterletaddonconfigs OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 16706)
-- Name: machinepools; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.machinepools (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.machinepools OWNER TO postgres;

--
-- TOC entry 207 (class 1259 OID 16419)
-- Name: placementbindings; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.placementbindings (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.placementbindings OWNER TO postgres;

--
-- TOC entry 206 (class 1259 OID 16405)
-- Name: placementrules; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.placementrules (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.placementrules OWNER TO postgres;

--
-- TOC entry 205 (class 1259 OID 16391)
-- Name: policies; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.policies (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.policies OWNER TO postgres;

--
-- TOC entry 231 (class 1259 OID 16680)
-- Name: secrets; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.secrets (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.secrets OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 16447)
-- Name: subscriptions; Type: TABLE; Schema: spec; Owner: postgres
--

CREATE TABLE spec.subscriptions (
    id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE spec.subscriptions OWNER TO postgres;

--
-- TOC entry 3046 (class 2606 OID 16485)
-- Name: applications applications_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.applications
    ADD CONSTRAINT applications_pkey PRIMARY KEY (id);


--
-- TOC entry 3044 (class 2606 OID 16471)
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- TOC entry 3048 (class 2606 OID 16499)
-- Name: clusterdeployments clusterdeployments_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.clusterdeployments
    ADD CONSTRAINT clusterdeployments_pkey PRIMARY KEY (id);


--
-- TOC entry 3040 (class 2606 OID 16443)
-- Name: configs configs_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.configs
    ADD CONSTRAINT configs_pkey PRIMARY KEY (id);


--
-- TOC entry 3054 (class 2606 OID 16730)
-- Name: klusterletaddonconfigs klusterletaddonconfigs_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.klusterletaddonconfigs
    ADD CONSTRAINT klusterletaddonconfigs_pkey PRIMARY KEY (id);


--
-- TOC entry 3052 (class 2606 OID 16716)
-- Name: machinepools machinepools_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.machinepools
    ADD CONSTRAINT machinepools_pkey PRIMARY KEY (id);


--
-- TOC entry 3038 (class 2606 OID 16429)
-- Name: placementbindings placementbindings_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.placementbindings
    ADD CONSTRAINT placementbindings_pkey PRIMARY KEY (id);


--
-- TOC entry 3036 (class 2606 OID 16415)
-- Name: placementrules placementrules_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.placementrules
    ADD CONSTRAINT placementrules_pkey PRIMARY KEY (id);


--
-- TOC entry 3034 (class 2606 OID 16401)
-- Name: policies policies_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.policies
    ADD CONSTRAINT policies_pkey PRIMARY KEY (id);


--
-- TOC entry 3050 (class 2606 OID 16690)
-- Name: secrets secrets_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.secrets
    ADD CONSTRAINT secrets_pkey PRIMARY KEY (id);


--
-- TOC entry 3042 (class 2606 OID 16457)
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: spec; Owner: postgres
--

ALTER TABLE ONLY spec.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- TOC entry 3067 (class 2620 OID 16488)
-- Name: applications move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.applications FOR EACH ROW EXECUTE FUNCTION public.move_applications_to_history();


--
-- TOC entry 3065 (class 2620 OID 16474)
-- Name: channels move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.channels FOR EACH ROW EXECUTE FUNCTION public.move_channels_to_history();


--
-- TOC entry 3069 (class 2620 OID 16502)
-- Name: clusterdeployments move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.clusterdeployments FOR EACH ROW EXECUTE FUNCTION public.move_clusterdeployments_to_history();


--
-- TOC entry 3061 (class 2620 OID 16446)
-- Name: configs move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.configs FOR EACH ROW EXECUTE FUNCTION public.move_configs_to_history();


--
-- TOC entry 3075 (class 2620 OID 16733)
-- Name: klusterletaddonconfigs move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.klusterletaddonconfigs FOR EACH ROW EXECUTE FUNCTION public.move_klusterletaddonconfigs_to_history();


--
-- TOC entry 3073 (class 2620 OID 16719)
-- Name: machinepools move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.machinepools FOR EACH ROW EXECUTE FUNCTION public.move_machinepools_to_history();


--
-- TOC entry 3059 (class 2620 OID 16432)
-- Name: placementbindings move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.placementbindings FOR EACH ROW EXECUTE FUNCTION public.move_placementbindings_to_history();


--
-- TOC entry 3057 (class 2620 OID 16418)
-- Name: placementrules move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.placementrules FOR EACH ROW EXECUTE FUNCTION public.move_placementrules_to_history();


--
-- TOC entry 3055 (class 2620 OID 16404)
-- Name: policies move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.policies FOR EACH ROW EXECUTE FUNCTION public.move_policies_to_history();


--
-- TOC entry 3071 (class 2620 OID 16693)
-- Name: secrets move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.secrets FOR EACH ROW EXECUTE FUNCTION public.move_secrets_to_history();


--
-- TOC entry 3063 (class 2620 OID 16460)
-- Name: subscriptions move_to_history; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER move_to_history BEFORE INSERT ON spec.subscriptions FOR EACH ROW EXECUTE FUNCTION public.move_subscriptions_to_history();


--
-- TOC entry 3068 (class 2620 OID 16486)
-- Name: applications set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.applications FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3066 (class 2620 OID 16472)
-- Name: channels set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.channels FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3070 (class 2620 OID 16500)
-- Name: clusterdeployments set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.clusterdeployments FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3062 (class 2620 OID 16444)
-- Name: configs set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.configs FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3076 (class 2620 OID 16731)
-- Name: klusterletaddonconfigs set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.klusterletaddonconfigs FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3074 (class 2620 OID 16717)
-- Name: machinepools set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.machinepools FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3060 (class 2620 OID 16430)
-- Name: placementbindings set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.placementbindings FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3058 (class 2620 OID 16416)
-- Name: placementrules set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.placementrules FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3056 (class 2620 OID 16402)
-- Name: policies set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.policies FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3072 (class 2620 OID 16691)
-- Name: secrets set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.secrets FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3064 (class 2620 OID 16458)
-- Name: subscriptions set_timestamp; Type: TRIGGER; Schema: spec; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON spec.subscriptions FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3212 (class 0 OID 0)
-- Dependencies: 10
-- Name: SCHEMA spec; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA spec TO hoh_process_user;
GRANT USAGE ON SCHEMA spec TO transport_bridge_user;


--
-- TOC entry 3213 (class 0 OID 0)
-- Dependencies: 211
-- Name: TABLE applications; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.applications TO hoh_process_user;
GRANT SELECT ON TABLE spec.applications TO transport_bridge_user;


--
-- TOC entry 3214 (class 0 OID 0)
-- Dependencies: 210
-- Name: TABLE channels; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.channels TO hoh_process_user;
GRANT SELECT ON TABLE spec.channels TO transport_bridge_user;


--
-- TOC entry 3215 (class 0 OID 0)
-- Dependencies: 212
-- Name: TABLE clusterdeployments; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.clusterdeployments TO hoh_process_user;
GRANT SELECT ON TABLE spec.clusterdeployments TO transport_bridge_user;


--
-- TOC entry 3216 (class 0 OID 0)
-- Dependencies: 208
-- Name: TABLE configs; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.configs TO hoh_process_user;
GRANT SELECT ON TABLE spec.configs TO transport_bridge_user;


--
-- TOC entry 3217 (class 0 OID 0)
-- Dependencies: 234
-- Name: TABLE klusterletaddonconfigs; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.klusterletaddonconfigs TO hoh_process_user;
GRANT SELECT ON TABLE spec.klusterletaddonconfigs TO transport_bridge_user;


--
-- TOC entry 3218 (class 0 OID 0)
-- Dependencies: 233
-- Name: TABLE machinepools; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.machinepools TO hoh_process_user;
GRANT SELECT ON TABLE spec.machinepools TO transport_bridge_user;


--
-- TOC entry 3219 (class 0 OID 0)
-- Dependencies: 207
-- Name: TABLE placementbindings; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.placementbindings TO hoh_process_user;
GRANT SELECT ON TABLE spec.placementbindings TO transport_bridge_user;


--
-- TOC entry 3220 (class 0 OID 0)
-- Dependencies: 206
-- Name: TABLE placementrules; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.placementrules TO hoh_process_user;
GRANT SELECT ON TABLE spec.placementrules TO transport_bridge_user;


--
-- TOC entry 3221 (class 0 OID 0)
-- Dependencies: 205
-- Name: TABLE policies; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.policies TO hoh_process_user;
GRANT SELECT ON TABLE spec.policies TO transport_bridge_user;


--
-- TOC entry 3222 (class 0 OID 0)
-- Dependencies: 231
-- Name: TABLE secrets; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.secrets TO hoh_process_user;
GRANT SELECT ON TABLE spec.secrets TO transport_bridge_user;


--
-- TOC entry 3223 (class 0 OID 0)
-- Dependencies: 209
-- Name: TABLE subscriptions; Type: ACL; Schema: spec; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE spec.subscriptions TO hoh_process_user;
GRANT SELECT ON TABLE spec.subscriptions TO transport_bridge_user;


-- Completed on 2021-10-21 14:06:27 EDT

--
-- PostgreSQL database dump complete
--


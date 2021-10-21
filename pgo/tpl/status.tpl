--
-- PostgreSQL database dump
--

-- Dumped from database version 13.4
-- Dumped by pg_dump version 14.0

-- Started on 2021-10-21 11:52:11 EDT

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
-- TOC entry 5 (class 2615 OID 16386)
-- Name: status; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA status;


ALTER SCHEMA status OWNER TO postgres;

--
-- TOC entry 734 (class 1247 OID 16634)
-- Name: compliance_type; Type: TYPE; Schema: status; Owner: postgres
--

CREATE TYPE status.compliance_type AS ENUM (
    'compliant',
    'non_compliant',
    'unknown'
);


ALTER TYPE status.compliance_type OWNER TO postgres;

--
-- TOC entry 727 (class 1247 OID 16620)
-- Name: error_type; Type: TYPE; Schema: status; Owner: postgres
--

CREATE TYPE status.error_type AS ENUM (
    'disconnected',
    'none'
);


ALTER TYPE status.error_type OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 226 (class 1259 OID 16651)
-- Name: aggregated_compliance; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.aggregated_compliance (
    id uuid,
    leaf_hub_name character varying(63) NOT NULL,
    applied_clusters integer NOT NULL,
    non_compliant_clusters integer NOT NULL
);


ALTER TABLE status.aggregated_compliance OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 16668)
-- Name: clusterdeployments; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.clusterdeployments (
    id uuid NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE status.clusterdeployments OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 16641)
-- Name: compliance; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.compliance (
    id uuid NOT NULL,
    cluster_name character varying(63) NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    error status.error_type NOT NULL,
    compliance status.compliance_type NOT NULL
);


ALTER TABLE status.compliance OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 16770)
-- Name: klusterletaddonconfigs; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.klusterletaddonconfigs (
    id uuid NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE status.klusterletaddonconfigs OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 16654)
-- Name: leaf_hub_heartbeats; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.leaf_hub_heartbeats (
    name character varying(63) NOT NULL,
    last_timestamp timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE status.leaf_hub_heartbeats OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 16762)
-- Name: machinepools; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.machinepools (
    id uuid NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE status.machinepools OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 16625)
-- Name: managed_clusters; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.managed_clusters (
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL,
    error status.error_type NOT NULL
);


ALTER TABLE status.managed_clusters OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 16674)
-- Name: secrets; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.secrets (
    id uuid NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE status.secrets OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 16660)
-- Name: subscriptions; Type: TABLE; Schema: status; Owner: postgres
--

CREATE TABLE status.subscriptions (
    id uuid NOT NULL,
    leaf_hub_name character varying(63) NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE status.subscriptions OWNER TO postgres;

--
-- TOC entry 3148 (class 0 OID 16651)
-- Dependencies: 226
-- Data for Name: aggregated_compliance; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.aggregated_compliance (id, leaf_hub_name, applied_clusters, non_compliant_clusters) FROM stdin;
\.


--
-- TOC entry 3151 (class 0 OID 16668)
-- Dependencies: 229
-- Data for Name: clusterdeployments; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.clusterdeployments (id, leaf_hub_name, payload) FROM stdin;
\.


--
-- TOC entry 3147 (class 0 OID 16641)
-- Dependencies: 224
-- Data for Name: compliance; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.compliance (id, cluster_name, leaf_hub_name, error, compliance) FROM stdin;
\.


--
-- TOC entry 3154 (class 0 OID 16770)
-- Dependencies: 238
-- Data for Name: klusterletaddonconfigs; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.klusterletaddonconfigs (id, leaf_hub_name, payload) FROM stdin;
\.


--
-- TOC entry 3149 (class 0 OID 16654)
-- Dependencies: 227
-- Data for Name: leaf_hub_heartbeats; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.leaf_hub_heartbeats (name, last_timestamp) FROM stdin;
\.


--
-- TOC entry 3153 (class 0 OID 16762)
-- Dependencies: 237
-- Data for Name: machinepools; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.machinepools (id, leaf_hub_name, payload) FROM stdin;
\.


--
-- TOC entry 3146 (class 0 OID 16625)
-- Dependencies: 223
-- Data for Name: managed_clusters; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.managed_clusters (leaf_hub_name, payload, error) FROM stdin;
hub1	{"kind": "ManagedCluster", "spec": {"hubAcceptsClient": true, "leaseDurationSeconds": 60}, "status": {"version": {}, "capacity": {"core_worker": "0", "socket_worker": "0"}, "conditions": [{"type": "HubAcceptedManagedCluster", "reason": "HubClusterAdminAccepted", "status": "True", "message": "Accepted by hub cluster admin", "lastTransitionTime": "2021-10-06T15:28:36Z"}, {"type": "ManagedClusterConditionAvailable", "reason": "ManagedClusterLeaseUpdateStopped", "status": "Unknown", "message": "Registration agent stopped updating its lease.", "lastTransitionTime": "2021-10-06T15:38:37Z"}]}, "metadata": {"uid": "e107a885-2af6-4b99-8dff-09282ae5ee57", "name": "cluster3", "labels": {"name": "cluster3", "vendor": "Kind", "feature.open-cluster-management.io/addon-work-manager": "unreachable", "feature.open-cluster-management.io/addon-search-collector": "unreachable", "feature.open-cluster-management.io/addon-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-application-manager": "unreachable", "feature.open-cluster-management.io/addon-iam-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-cert-policy-controller": "unreachable"}, "annotations": {"open-cluster-management/created-via": "other"}, "resourceVersion": "8689023", "creationTimestamp": "2021-10-06T15:28:36Z"}, "apiVersion": "cluster.open-cluster-management.io/v1"}	none
hub1	{"kind": "ManagedCluster", "spec": {"hubAcceptsClient": true, "leaseDurationSeconds": 60}, "status": {"version": {}, "capacity": {"core_worker": "0", "socket_worker": "0"}, "conditions": [{"type": "HubAcceptedManagedCluster", "reason": "HubClusterAdminAccepted", "status": "True", "message": "Accepted by hub cluster admin", "lastTransitionTime": "2021-10-06T15:28:10Z"}, {"type": "ManagedClusterConditionAvailable", "reason": "ManagedClusterLeaseUpdateStopped", "status": "Unknown", "message": "Registration agent stopped updating its lease.", "lastTransitionTime": "2021-10-06T15:38:10Z"}]}, "metadata": {"uid": "51e97333-01b8-416e-b2aa-4268f83868d7", "name": "cluster1", "labels": {"name": "cluster1", "vendor": "Kind", "feature.open-cluster-management.io/addon-work-manager": "unreachable", "feature.open-cluster-management.io/addon-search-collector": "unreachable", "feature.open-cluster-management.io/addon-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-application-manager": "unreachable", "feature.open-cluster-management.io/addon-iam-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-cert-policy-controller": "unreachable"}, "annotations": {"open-cluster-management/created-via": "other"}, "resourceVersion": "8689019", "creationTimestamp": "2021-10-06T15:28:09Z"}, "apiVersion": "cluster.open-cluster-management.io/v1"}	none
hub1	{"kind": "ManagedCluster", "spec": {"hubAcceptsClient": true, "leaseDurationSeconds": 60}, "status": {"version": {}, "capacity": {"core_worker": "0", "socket_worker": "0"}, "conditions": [{"type": "HubAcceptedManagedCluster", "reason": "HubClusterAdminAccepted", "status": "True", "message": "Accepted by hub cluster admin", "lastTransitionTime": "2021-10-06T15:28:22Z"}, {"type": "ManagedClusterConditionAvailable", "reason": "ManagedClusterLeaseUpdateStopped", "status": "Unknown", "message": "Registration agent stopped updating its lease.", "lastTransitionTime": "2021-10-06T15:38:22Z"}]}, "metadata": {"uid": "71abcc70-96fe-4541-a9de-95b5a32556d3", "name": "cluster2", "labels": {"name": "cluster2", "vendor": "Kind", "feature.open-cluster-management.io/addon-work-manager": "unreachable", "feature.open-cluster-management.io/addon-search-collector": "unreachable", "feature.open-cluster-management.io/addon-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-application-manager": "unreachable", "feature.open-cluster-management.io/addon-iam-policy-controller": "unreachable", "feature.open-cluster-management.io/addon-cert-policy-controller": "unreachable"}, "annotations": {"open-cluster-management/created-via": "other"}, "resourceVersion": "8689020", "creationTimestamp": "2021-10-06T15:28:22Z"}, "apiVersion": "cluster.open-cluster-management.io/v1"}	none
hub1	{"kind": "ManagedCluster", "spec": {"hubAcceptsClient": true, "leaseDurationSeconds": 60}, "status": {"version": {"kubernetes": "v1.19.1"}, "capacity": {"cpu": "8", "pods": "110", "memory": "32886516Ki", "core_worker": "0", "hugepages-1Gi": "0", "hugepages-2Mi": "0", "socket_worker": "0", "ephemeral-storage": "101583780Ki"}, "conditions": [{"type": "HubAcceptedManagedCluster", "reason": "HubClusterAdminAccepted", "status": "True", "message": "Accepted by hub cluster admin", "lastTransitionTime": "2021-10-06T15:28:49Z"}, {"type": "ManagedClusterJoined", "reason": "ManagedClusterJoined", "status": "True", "message": "Managed cluster joined", "lastTransitionTime": "2021-10-06T15:29:17Z"}, {"type": "ManagedClusterConditionAvailable", "reason": "ManagedClusterAvailable", "status": "True", "message": "Managed cluster is available", "lastTransitionTime": "2021-10-19T14:31:32Z"}], "allocatable": {"cpu": "8", "pods": "110", "memory": "32886516Ki", "hugepages-1Gi": "0", "hugepages-2Mi": "0", "ephemeral-storage": "101583780Ki"}, "clusterClaims": [{"name": "id.k8s.io", "value": "cluster4"}, {"name": "kubeversion.open-cluster-management.io", "value": "v1.19.1"}, {"name": "platform.open-cluster-management.io", "value": "Other"}, {"name": "product.open-cluster-management.io", "value": "Other"}]}, "metadata": {"uid": "ed7a9c58-8883-462e-826c-c1d11858f36d", "name": "cluster4", "labels": {"name": "cluster4", "vendor": "Kind", "feature.open-cluster-management.io/addon-work-manager": "available", "feature.open-cluster-management.io/addon-search-collector": "available", "feature.open-cluster-management.io/addon-policy-controller": "available", "feature.open-cluster-management.io/addon-application-manager": "available", "feature.open-cluster-management.io/addon-iam-policy-controller": "available", "feature.open-cluster-management.io/addon-cert-policy-controller": "available"}, "annotations": {"open-cluster-management/created-via": "other"}, "resourceVersion": "12525916", "creationTimestamp": "2021-10-06T15:28:48Z"}, "apiVersion": "cluster.open-cluster-management.io/v1"}	none
hub1	{"kind": "ManagedCluster", "spec": {"hubAcceptsClient": true, "leaseDurationSeconds": 60}, "status": {"version": {"kubernetes": "v1.19.1"}, "capacity": {"cpu": "8", "pods": "110", "memory": "32886516Ki", "core_worker": "0", "hugepages-1Gi": "0", "hugepages-2Mi": "0", "socket_worker": "0", "ephemeral-storage": "101583780Ki"}, "conditions": [{"type": "HubAcceptedManagedCluster", "reason": "HubClusterAdminAccepted", "status": "True", "message": "Accepted by hub cluster admin", "lastTransitionTime": "2021-10-06T15:27:56Z"}, {"type": "ManagedClusterJoined", "reason": "ManagedClusterJoined", "status": "True", "message": "Managed cluster joined", "lastTransitionTime": "2021-10-06T15:28:19Z"}, {"type": "ManagedClusterConditionAvailable", "reason": "ManagedClusterAvailable", "status": "True", "message": "Managed cluster is available", "lastTransitionTime": "2021-10-19T14:31:19Z"}], "allocatable": {"cpu": "8", "pods": "110", "memory": "32886516Ki", "hugepages-1Gi": "0", "hugepages-2Mi": "0", "ephemeral-storage": "101583780Ki"}, "clusterClaims": [{"name": "id.k8s.io", "value": "cluster0"}, {"name": "kubeversion.open-cluster-management.io", "value": "v1.19.1"}, {"name": "platform.open-cluster-management.io", "value": "Other"}, {"name": "product.open-cluster-management.io", "value": "Other"}]}, "metadata": {"uid": "600ce8a7-a1f0-4d96-b7f2-876992db925d", "name": "cluster0", "labels": {"name": "cluster0", "vendor": "Kind", "feature.open-cluster-management.io/addon-work-manager": "available", "feature.open-cluster-management.io/addon-cluster-proxy": "unreachable", "feature.open-cluster-management.io/addon-search-collector": "available", "feature.open-cluster-management.io/addon-policy-controller": "available", "feature.open-cluster-management.io/addon-application-manager": "available", "feature.open-cluster-management.io/addon-iam-policy-controller": "available", "feature.open-cluster-management.io/addon-cert-policy-controller": "available"}, "annotations": {"open-cluster-management/created-via": "other"}, "resourceVersion": "12524907", "creationTimestamp": "2021-10-06T15:27:56Z"}, "apiVersion": "cluster.open-cluster-management.io/v1"}	none
\.


--
-- TOC entry 3152 (class 0 OID 16674)
-- Dependencies: 230
-- Data for Name: secrets; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.secrets (id, leaf_hub_name, payload) FROM stdin;
\.


--
-- TOC entry 3150 (class 0 OID 16660)
-- Dependencies: 228
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: status; Owner: postgres
--

COPY status.subscriptions (id, leaf_hub_name, payload) FROM stdin;
\.


--
-- TOC entry 3007 (class 2606 OID 16659)
-- Name: leaf_hub_heartbeats leaf_hub_heartbeats_pkey; Type: CONSTRAINT; Schema: status; Owner: postgres
--

ALTER TABLE ONLY status.leaf_hub_heartbeats
    ADD CONSTRAINT leaf_hub_heartbeats_pkey PRIMARY KEY (name);


--
-- TOC entry 3008 (class 1259 OID 16758)
-- Name: clusterdeployments_metadata_name_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX clusterdeployments_metadata_name_idx ON status.clusterdeployments USING btree ((((payload -> 'metadata'::text) ->> 'name'::text)));


--
-- TOC entry 3009 (class 1259 OID 16759)
-- Name: clusterdeployments_metadata_namespace_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX clusterdeployments_metadata_namespace_idx ON status.clusterdeployments USING btree ((((payload -> 'metadata'::text) ->> 'namespace'::text)));


--
-- TOC entry 3002 (class 1259 OID 16645)
-- Name: compliance_id_non_compliant_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE INDEX compliance_id_non_compliant_idx ON status.compliance USING btree (id, compliance) WHERE (compliance = 'non_compliant'::status.compliance_type);


--
-- TOC entry 3003 (class 1259 OID 16646)
-- Name: compliance_leaf_hub_cluster_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE INDEX compliance_leaf_hub_cluster_idx ON status.compliance USING btree (leaf_hub_name, cluster_name);


--
-- TOC entry 3004 (class 1259 OID 16644)
-- Name: compliance_leaf_hub_non_compliant_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE INDEX compliance_leaf_hub_non_compliant_idx ON status.compliance USING btree (leaf_hub_name, compliance) WHERE (compliance = 'non_compliant'::status.compliance_type);


--
-- TOC entry 3005 (class 1259 OID 16647)
-- Name: compliance_policy_leaf_hub_cluster_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX compliance_policy_leaf_hub_cluster_idx ON status.compliance USING btree (id, leaf_hub_name, cluster_name);


--
-- TOC entry 3014 (class 1259 OID 16776)
-- Name: klusterletaddonconfigs_metadata_name_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX klusterletaddonconfigs_metadata_name_idx ON status.klusterletaddonconfigs USING btree ((((payload -> 'metadata'::text) ->> 'name'::text)));


--
-- TOC entry 3015 (class 1259 OID 16777)
-- Name: klusterletaddonconfigs_metadata_namespace_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX klusterletaddonconfigs_metadata_namespace_idx ON status.klusterletaddonconfigs USING btree ((((payload -> 'metadata'::text) ->> 'namespace'::text)));


--
-- TOC entry 3012 (class 1259 OID 16768)
-- Name: machinepools_metadata_name_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX machinepools_metadata_name_idx ON status.machinepools USING btree ((((payload -> 'metadata'::text) ->> 'name'::text)));


--
-- TOC entry 3013 (class 1259 OID 16769)
-- Name: machinepools_metadata_namespace_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX machinepools_metadata_namespace_idx ON status.machinepools USING btree ((((payload -> 'metadata'::text) ->> 'namespace'::text)));


--
-- TOC entry 3001 (class 1259 OID 16632)
-- Name: managed_clusters_metadata_name_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX managed_clusters_metadata_name_idx ON status.managed_clusters USING btree ((((payload -> 'metadata'::text) ->> 'name'::text)));


--
-- TOC entry 3010 (class 1259 OID 16760)
-- Name: secrets_metadata_name_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX secrets_metadata_name_idx ON status.secrets USING btree ((((payload -> 'metadata'::text) ->> 'name'::text)));


--
-- TOC entry 3011 (class 1259 OID 16761)
-- Name: secrets_metadata_namespace_idx; Type: INDEX; Schema: status; Owner: postgres
--

CREATE UNIQUE INDEX secrets_metadata_namespace_idx ON status.secrets USING btree ((((payload -> 'metadata'::text) ->> 'namespace'::text)));


--
-- TOC entry 3160 (class 0 OID 0)
-- Dependencies: 5
-- Name: SCHEMA status; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA status TO hoh_process_user;
GRANT USAGE ON SCHEMA status TO transport_bridge_user;


--
-- TOC entry 3161 (class 0 OID 0)
-- Dependencies: 226
-- Name: TABLE aggregated_compliance; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.aggregated_compliance TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.aggregated_compliance TO transport_bridge_user;


--
-- TOC entry 3162 (class 0 OID 0)
-- Dependencies: 229
-- Name: TABLE clusterdeployments; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.clusterdeployments TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.clusterdeployments TO transport_bridge_user;


--
-- TOC entry 3163 (class 0 OID 0)
-- Dependencies: 224
-- Name: TABLE compliance; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.compliance TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.compliance TO transport_bridge_user;


--
-- TOC entry 3164 (class 0 OID 0)
-- Dependencies: 238
-- Name: TABLE klusterletaddonconfigs; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.klusterletaddonconfigs TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.klusterletaddonconfigs TO transport_bridge_user;


--
-- TOC entry 3165 (class 0 OID 0)
-- Dependencies: 227
-- Name: TABLE leaf_hub_heartbeats; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.leaf_hub_heartbeats TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.leaf_hub_heartbeats TO transport_bridge_user;


--
-- TOC entry 3166 (class 0 OID 0)
-- Dependencies: 237
-- Name: TABLE machinepools; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.machinepools TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.machinepools TO transport_bridge_user;


--
-- TOC entry 3167 (class 0 OID 0)
-- Dependencies: 223
-- Name: TABLE managed_clusters; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.managed_clusters TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.managed_clusters TO transport_bridge_user;


--
-- TOC entry 3168 (class 0 OID 0)
-- Dependencies: 230
-- Name: TABLE secrets; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.secrets TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.secrets TO transport_bridge_user;


--
-- TOC entry 3169 (class 0 OID 0)
-- Dependencies: 228
-- Name: TABLE subscriptions; Type: ACL; Schema: status; Owner: postgres
--

GRANT SELECT ON TABLE status.subscriptions TO hoh_process_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE status.subscriptions TO transport_bridge_user;


-- Completed on 2021-10-21 11:52:16 EDT

--
-- PostgreSQL database dump complete
--


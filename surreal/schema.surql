-- Headlamp Agent SurrealDB schema
-- SurrealDB is the operational graph, artifact, workflow, live-event, and current-state store.

DEFINE NAMESPACE IF NOT EXISTS headlamp;
DEFINE DATABASE IF NOT EXISTS headlamp;

-- Core artifact/state tables
DEFINE TABLE IF NOT EXISTS hosts SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS kind ON hosts TYPE string ASSERT $value = "HostServerArtifact";
DEFINE FIELD IF NOT EXISTS name ON hosts TYPE string;
DEFINE FIELD IF NOT EXISTS type ON hosts TYPE string ASSERT $value = "host-server";
DEFINE FIELD IF NOT EXISTS environment ON hosts TYPE option<string>;
DEFINE FIELD IF NOT EXISTS provider ON hosts TYPE option<string>;
DEFINE FIELD IF NOT EXISTS region ON hosts TYPE option<string>;
DEFINE FIELD IF NOT EXISTS roles ON hosts TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS status ON hosts TYPE object;
DEFINE FIELD IF NOT EXISTS labels ON hosts TYPE option<object>;
DEFINE FIELD IF NOT EXISTS created_at ON hosts TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON hosts TYPE datetime DEFAULT time::now();

DEFINE TABLE IF NOT EXISTS agents SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS kind ON agents TYPE string ASSERT $value = "AgentArtifact";
DEFINE FIELD IF NOT EXISTS name ON agents TYPE string;
DEFINE FIELD IF NOT EXISTS type ON agents TYPE string ASSERT $value = "agent";
DEFINE FIELD IF NOT EXISTS host_id ON agents TYPE record<hosts>;
DEFINE FIELD IF NOT EXISTS mode ON agents TYPE option<string>;
DEFINE FIELD IF NOT EXISTS version ON agents TYPE option<string>;
DEFINE FIELD IF NOT EXISTS capabilities ON agents TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS status ON agents TYPE object;
DEFINE FIELD IF NOT EXISTS created_at ON agents TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON agents TYPE datetime DEFAULT time::now();

DEFINE TABLE IF NOT EXISTS clusters SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS kind ON clusters TYPE string ASSERT $value = "ClusterInstanceArtifact";
DEFINE FIELD IF NOT EXISTS name ON clusters TYPE string;
DEFINE FIELD IF NOT EXISTS type ON clusters TYPE string ASSERT $value = "cluster-instance";
DEFINE FIELD IF NOT EXISTS distribution ON clusters TYPE string;
DEFINE FIELD IF NOT EXISTS environment ON clusters TYPE option<string>;
DEFINE FIELD IF NOT EXISTS provider ON clusters TYPE option<string>;
DEFINE FIELD IF NOT EXISTS region ON clusters TYPE option<string>;
DEFINE FIELD IF NOT EXISTS endpoint ON clusters TYPE object;
DEFINE FIELD IF NOT EXISTS version ON clusters TYPE option<object>;
DEFINE FIELD IF NOT EXISTS topology ON clusters TYPE option<object>;
DEFINE FIELD IF NOT EXISTS network ON clusters TYPE option<object>;
DEFINE FIELD IF NOT EXISTS security ON clusters TYPE option<object>;
DEFINE FIELD IF NOT EXISTS observability ON clusters TYPE option<object>;
DEFINE FIELD IF NOT EXISTS status ON clusters TYPE object;
DEFINE FIELD IF NOT EXISTS created_at ON clusters TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON clusters TYPE datetime DEFAULT time::now();

-- Live/event tables
DEFINE TABLE IF NOT EXISTS agent_heartbeats SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS agent_id ON agent_heartbeats TYPE record<agents>;
DEFINE FIELD IF NOT EXISTS host_id ON agent_heartbeats TYPE record<hosts>;
DEFINE FIELD IF NOT EXISTS status ON agent_heartbeats TYPE string;
DEFINE FIELD IF NOT EXISTS version ON agent_heartbeats TYPE option<string>;
DEFINE FIELD IF NOT EXISTS capabilities ON agent_heartbeats TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS health ON agent_heartbeats TYPE option<object>;
DEFINE FIELD IF NOT EXISTS observed_at ON agent_heartbeats TYPE datetime DEFAULT time::now();
DEFINE INDEX IF NOT EXISTS agent_heartbeats_agent_time ON agent_heartbeats FIELDS agent_id, observed_at;

DEFINE TABLE IF NOT EXISTS health_events SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS subject_kind ON health_events TYPE string;
DEFINE FIELD IF NOT EXISTS subject_id ON health_events TYPE string;
DEFINE FIELD IF NOT EXISTS previous_status ON health_events TYPE option<string>;
DEFINE FIELD IF NOT EXISTS current_status ON health_events TYPE string;
DEFINE FIELD IF NOT EXISTS severity ON health_events TYPE string;
DEFINE FIELD IF NOT EXISTS reason ON health_events TYPE option<string>;
DEFINE FIELD IF NOT EXISTS created_at ON health_events TYPE datetime DEFAULT time::now();

DEFINE TABLE IF NOT EXISTS notification_events SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS severity ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS source ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS subject_kind ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS subject_id ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS title ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS message ON notification_events TYPE string;
DEFINE FIELD IF NOT EXISTS status ON notification_events TYPE string DEFAULT "open";
DEFINE FIELD IF NOT EXISTS recommended_actions ON notification_events TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS channels ON notification_events TYPE option<array<object>>;
DEFINE FIELD IF NOT EXISTS created_at ON notification_events TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS resolved_at ON notification_events TYPE option<datetime>;
DEFINE INDEX IF NOT EXISTS notification_subject_status ON notification_events FIELDS subject_id, status;

DEFINE TABLE IF NOT EXISTS incidents SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS severity ON incidents TYPE string;
DEFINE FIELD IF NOT EXISTS status ON incidents TYPE string DEFAULT "open";
DEFINE FIELD IF NOT EXISTS title ON incidents TYPE string;
DEFINE FIELD IF NOT EXISTS affected_artifacts ON incidents TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS detected_by ON incidents TYPE string;
DEFINE FIELD IF NOT EXISTS recovery_state ON incidents TYPE option<string>;
DEFINE FIELD IF NOT EXISTS next_action ON incidents TYPE option<string>;
DEFINE FIELD IF NOT EXISTS created_at ON incidents TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON incidents TYPE datetime DEFAULT time::now();
DEFINE INDEX IF NOT EXISTS incidents_status_severity ON incidents FIELDS status, severity;

DEFINE TABLE IF NOT EXISTS opa_decisions SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS policy_bundle ON opa_decisions TYPE string;
DEFINE FIELD IF NOT EXISTS subject_id ON opa_decisions TYPE string;
DEFINE FIELD IF NOT EXISTS action ON opa_decisions TYPE string;
DEFINE FIELD IF NOT EXISTS result ON opa_decisions TYPE string;
DEFINE FIELD IF NOT EXISTS reasons ON opa_decisions TYPE option<array<string>>;
DEFINE FIELD IF NOT EXISTS input ON opa_decisions TYPE object;
DEFINE FIELD IF NOT EXISTS created_at ON opa_decisions TYPE datetime DEFAULT time::now();

DEFINE TABLE IF NOT EXISTS audit_events SCHEMAFULL;
DEFINE FIELD IF NOT EXISTS actor ON audit_events TYPE string;
DEFINE FIELD IF NOT EXISTS channel ON audit_events TYPE string;
DEFINE FIELD IF NOT EXISTS subject_id ON audit_events TYPE string;
DEFINE FIELD IF NOT EXISTS action ON audit_events TYPE string;
DEFINE FIELD IF NOT EXISTS decision ON audit_events TYPE option<string>;
DEFINE FIELD IF NOT EXISTS result ON audit_events TYPE string;
DEFINE FIELD IF NOT EXISTS metadata ON audit_events TYPE option<object>;
DEFINE FIELD IF NOT EXISTS created_at ON audit_events TYPE datetime DEFAULT time::now();
DEFINE INDEX IF NOT EXISTS audit_subject_time ON audit_events FIELDS subject_id, created_at;

-- Graph relationship tables
DEFINE TABLE IF NOT EXISTS runs_on TYPE RELATION IN agents OUT hosts;
DEFINE TABLE IF NOT EXISTS belongs_to_cluster TYPE RELATION IN hosts OUT clusters;
DEFINE TABLE IF NOT EXISTS affected_by TYPE RELATION IN incidents OUT hosts | clusters | agents;
DEFINE TABLE IF NOT EXISTS governed_by TYPE RELATION IN clusters | agents | hosts OUT opa_decisions;

-- Useful live subscriptions for Headlamp Desktop / Cortex:
-- LIVE SELECT * FROM agent_heartbeats;
-- LIVE SELECT * FROM notification_events WHERE status = "open";
-- LIVE SELECT * FROM incidents WHERE status = "open";
-- LIVE SELECT * FROM health_events;

<template>
  <div class="card-title">
    <h2>{{ $t("settings.adminDashboard") }}</h2>
    <span class="refresh-badge" :class="{ refreshing }">
      <i class="material-icons">sync</i>
    </span>
  </div>

  <div class="card-content dashboard-content">
    <!-- Top stat cards -->
    <div class="stat-grid">
      <div class="stat-card">
        <i class="material-icons stat-icon">schedule</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashUptime") }}</p>
          <p class="stat-value">{{ formatUptime(stats.uptime) }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">info</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("general.version") }}</p>
          <p class="stat-value">{{ stats.version || "—" }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">memory</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashMemory") }}</p>
          <p class="stat-value">{{ formatBytes(stats.memAlloc) }}</p>
          <p class="stat-sub">{{ $t("settings.dashHeapSys") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ formatBytes(stats.heapSys) }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">dns</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashStack") }}</p>
          <p class="stat-value">{{ formatBytes(stats.stackInuse) }}</p>
          <p class="stat-sub">{{ $t("settings.dashMemSys") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ formatBytes(stats.memSys) }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">settings_ethernet</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashGoroutines") }}</p>
          <p class="stat-value">{{ stats.goroutines ?? "—" }}</p>
          <p class="stat-sub">{{ $t("settings.dashCPUs") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ stats.numCPU ?? "—" }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">recycling</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashGCPause") }}</p>
          <p class="stat-value">{{ stats.pauseTotalMs ?? 0 }}ms</p><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          <p class="stat-sub">GC runs<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ stats.numGC ?? "—" }}</p>
        </div>
      </div>
      <div class="stat-card">
        <i class="material-icons stat-icon">moving</i>
        <div class="stat-body">
          <p class="stat-label">{{ $t("settings.dashTotalAlloc") }}</p>
          <p class="stat-value">{{ formatBytes(stats.totalAlloc) }}</p>
          <p class="stat-sub">{{ $t("settings.dashLastGC") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ formatUptime(stats.lastGCSecs) }} ago</p>
        </div>
      </div>
    </div>

    <!-- Per-source cards -->
    <h3 class="section-heading">{{ $t("general.sources") }}</h3>
    <div v-if="Object.keys(sources).length === 0" class="empty-state">
      <i class="material-icons">storage</i>
      <p>{{ $t("general.loading", { suffix: "..." }) }}</p>
    </div>
    <div v-else class="source-list">
      <div v-for="(src, name) in sources" :key="name" class="source-row">

        <!-- Info tile -->
        <div class="source-card">
          <div class="source-header">
            <i class="material-icons">storage</i>
            <span class="source-name">{{ src.name || name }}</span>
            <span class="status-badge" :class="src.status">{{ src.status }}</span>
          </div>

          <!-- Disk usage bar -->
          <div class="disk-bar-wrap" v-if="src.total > 0">
            <div class="disk-bar-track">
              <div class="disk-bar-fill" :style="{ width: diskPct(src) + '%' }"
                :class="diskClass(src)"></div>
            </div>
            <span class="disk-label">{{ formatBytes(src.used) }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text --> / {{ formatBytes(src.total) }} ({{ diskPct(src) }}%)</span> <!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          </div>

          <div class="source-stats">
            <div class="source-stat">
              <i class="material-icons">insert_drive_file</i>
              <span>{{ src.numFiles?.toLocaleString() ?? 0 }} {{ $t("general.files") }}</span>
            </div>
            <div class="source-stat">
              <i class="material-icons">folder</i>
              <span>{{ src.numDirs?.toLocaleString() ?? 0 }} {{ $t("general.folders") }}</span>
            </div>
            <div class="source-stat" v-if="src.numDeleted > 0">
              <i class="material-icons">delete_outline</i>
              <span>{{ $t("settings.dashNumDeleted") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ src.numDeleted?.toLocaleString() ?? 0 }}</span>
            </div>
            <div class="source-stat">
              <i class="material-icons">timer</i>
              <span>{{ $t("index.lastScanned") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ lastScanned(src) }}</span>
            </div>
            <div class="source-stat">
              <i class="material-icons">speed</i>
              <span>{{ $t("index.quickScan") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ src.quickScanDurationSeconds ?? 0 }}s &nbsp;|&nbsp; {{ $t("index.fullScan") }}: {{ src.fullScanDurationSeconds ?? 0 }}s</span> <!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
            </div>
            <div class="source-stat" v-if="src.complexity > 0">
              <i class="material-icons">bar_chart</i>
              <span>{{ $t("settings.dashComplexity") }}<!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->: {{ complexityLabel(src.complexity) }}</span>
            </div>
          </div>
        </div>

        <!-- Donut chart tile -->
        <div class="source-card chart-card" v-if="src.total > 0">
          <p class="stat-label">Storage Usage</p><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          <div class="donut-wrap">
            <svg class="donut-svg" viewBox="0 0 140 140">
              <!-- Background track -->
              <circle class="donut-track" cx="70" cy="70" r="54" />
              <!-- Used arc -->
              <circle
                class="donut-arc"
                :class="diskClass(src)"
                cx="70" cy="70" r="54"
                :stroke-dasharray="`${diskPct(src) * 3.393} 339.3`"
                stroke-dashoffset="84.8"
              />
            </svg>
            <div class="donut-label">
              <span class="donut-pct">{{ diskPct(src) }}%</span><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
              <span class="donut-sub">used</span><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
            </div>
          </div>
          <div class="donut-legend">
            <span class="legend-dot" :class="diskClass(src)"></span>
            <span>{{ formatBytes(src.used) }} used</span><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          </div>
          <div class="donut-legend">
            <span class="legend-dot free"></span>
            <span>{{ formatBytes(src.total - src.used) }} free</span><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          </div>
          <div class="donut-legend">
            <span class="legend-dot total"></span>
            <span>{{ formatBytes(src.total) }} total</span><!-- eslint-disable-line @intlify/vue-i18n/no-raw-text -->
          </div>
        </div>

      </div>
    </div>

    <p v-if="error" class="error-msg">
      <i class="material-icons">error_outline</i> {{ error }}
    </p>
  </div>
</template>

<script>
import { getStats } from "@/api/admin";
import { fromNow } from "@/utils/moment";
import { state } from "@/store";

export default {
  name: "AdminDashboard",
  data() {
    return {
      stats: {},
      refreshing: false,
      error: null,
      interval: null,
    };
  },
  computed: {
    sources() {
      return this.stats.sources || {};
    },
  },
  async mounted() {
    await this.load();
    this.interval = setInterval(this.load, 30000);
  },
  beforeUnmount() {
    clearInterval(this.interval);
  },
  methods: {
    async load() {
      this.refreshing = true;
      this.error = null;
      try {
        this.stats = await getStats();
      } catch (e) {
        this.error = e?.message || "Failed to load stats";
      } finally {
        this.refreshing = false;
      }
    },
    formatUptime(seconds) {
      if (!seconds) return "—";
      const d = Math.floor(seconds / 86400);
      const h = Math.floor((seconds % 86400) / 3600);
      const m = Math.floor((seconds % 3600) / 60);
      const s = seconds % 60;
      if (d > 0) return `${d}d ${h}h ${m}m`;
      if (h > 0) return `${h}h ${m}m ${s}s`;
      if (m > 0) return `${m}m ${s}s`;
      return `${s}s`;
    },
    formatBytes(bytes) {
      if (!bytes) return "0 B";
      const units = ["B", "KB", "MB", "GB", "TB"];
      let i = 0;
      let v = bytes;
      while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
      return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
    },
    diskPct(src) {
      if (!src.total || src.total === 0) return 0;
      return Math.round((src.used / src.total) * 100);
    },
    diskClass(src) {
      const pct = this.diskPct(src);
      if (pct >= 90) return "critical";
      if (pct >= 75) return "warning";
      return "ok";
    },
    lastScanned(src) {
      if (!src.lastIndexedUnixTime || src.lastIndexedUnixTime === 0) return "—";
      return fromNow(src.lastIndexedUnixTime, state.user?.locale);
    },
    complexityLabel(c) {
      if (c === 0) return "Unknown";
      if (c === 1) return "Simple";
      if (c <= 6) return "Normal";
      if (c <= 9) return "Complex";
      return "Highly Complex";
    },
  },
};
</script>

<style scoped>
.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 1.5em;
}

/* Refresh spin indicator */
.refresh-badge {
  display: inline-flex;
  align-items: center;
  margin-left: 0.75em;
  opacity: 0.3;
  transition: opacity 0.3s;
}
.refresh-badge.refreshing {
  opacity: 1;
}
.refresh-badge i {
  font-size: 1.1em;
  animation: none;
}
.refresh-badge.refreshing i {
  animation: spin 1s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Stat cards row */
.stat-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1em;
  justify-content: center;
}

.stat-card {
  display: flex;
  align-items: flex-start;
  gap: 0.75em;
  padding: 1.3em 1em;
  border-radius: 10px;
  background: var(--surfaceSecondary, #f3f4f6);
  min-height: 5.5em;
  flex: 1 1 180px;
  max-width: calc(25% - 0.75em);
}

@media (max-width: 700px) {
  .stat-card {
    max-width: calc(50% - 0.5em);
  }
}

.stat-icon {
  font-size: 1.6em;
  color: var(--primaryColor);
  margin-top: 0.1em;
  flex-shrink: 0;
}

.stat-body {
  min-width: 0;
}

.stat-label {
  font-size: 0.75em;
  color: var(--textSecondary, #888);
  margin: 0 0 0.2em;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 1.1em;
  font-weight: 600;
  margin: 0;
}

.stat-sub {
  font-size: 0.75em;
  color: var(--textSecondary, #888);
  margin: 0.2em 0 0;
}

/* Section heading */
.section-heading {
  font-size: 0.85em;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--textSecondary, #888);
  margin: 0;
}

/* Source list */
.source-list {
  display: flex;
  flex-direction: column;
  gap: 1em;
}

.source-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1em;
  align-items: stretch;
}

@media (max-width: 700px) {
  .source-row {
    grid-template-columns: 1fr;
  }
}

.source-card {
  padding: 2em;
  border-radius: 10px;
  background: var(--surfaceSecondary, #f3f4f6);
  display: flex;
  flex-direction: column;
  gap: 0.75em;
  min-height: 18em;
}

/* Chart tile */
.chart-card {
  align-items: center;
  justify-content: center;
}

.donut-wrap {
  position: relative;
  width: 140px;
  height: 140px;
}

.donut-svg {
  width: 140px;
  height: 140px;
  transform: rotate(-90deg);
}

.donut-track {
  fill: none;
  stroke: rgba(0,0,0,0.08);
  stroke-width: 16;
}

.donut-arc {
  fill: none;
  stroke-width: 16;
  stroke-linecap: round;
  transition: stroke-dasharray 0.5s ease;
}
.donut-arc.ok    { stroke: var(--primaryColor); }
.donut-arc.warning  { stroke: #f59e0b; }
.donut-arc.critical { stroke: #ef4444; }

.donut-label {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.donut-pct {
  font-size: 1.5em;
  font-weight: 700;
  line-height: 1;
}

.donut-sub {
  font-size: 0.7em;
  color: var(--textSecondary, #888);
  margin-top: 0.2em;
}

.donut-legend {
  display: flex;
  align-items: center;
  gap: 0.5em;
  font-size: 0.8em;
  color: var(--textSecondary, #666);
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.legend-dot.ok       { background: var(--primaryColor); }
.legend-dot.warning  { background: #f59e0b; }
.legend-dot.critical { background: #ef4444; }
.legend-dot.free     { background: rgba(0,0,0,0.12); }
.legend-dot.total    { background: transparent; border: 2px solid rgba(0,0,0,0.2); }

.source-header {
  display: flex;
  align-items: center;
  gap: 0.5em;
}

.source-header i {
  color: var(--primaryColor);
}

.source-name {
  font-weight: 600;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Status badge */
.status-badge {
  font-size: 0.7em;
  padding: 0.2em 0.6em;
  border-radius: 99px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  background: #ccc;
  color: #333;
}
.status-badge.ready { background: #d1fae5; color: #065f46; }
.status-badge.indexing { background: #fef3c7; color: #92400e; }
.status-badge.error, .status-badge.unavailable { background: #fee2e2; color: #991b1b; }

/* Disk bar */
.disk-bar-wrap {
  display: flex;
  flex-direction: column;
  gap: 0.3em;
}

.disk-bar-track {
  height: 6px;
  border-radius: 99px;
  background: rgba(0,0,0,0.1);
  overflow: hidden;
}

.disk-bar-fill {
  height: 100%;
  border-radius: 99px;
  transition: width 0.4s ease;
}
.disk-bar-fill.ok { background: var(--primaryColor); }
.disk-bar-fill.warning { background: #f59e0b; }
.disk-bar-fill.critical { background: #ef4444; }

.disk-label {
  font-size: 0.75em;
  color: var(--textSecondary, #888);
}

/* Source stat rows */
.source-stats {
  display: flex;
  flex-direction: column;
  gap: 0.3em;
}

.source-stat {
  display: flex;
  align-items: center;
  gap: 0.4em;
  font-size: 0.82em;
  color: var(--textSecondary, #666);
}

.source-stat i {
  font-size: 1em;
  color: var(--primaryColor);
  opacity: 0.7;
}

/* Empty / error */
.empty-state {
  display: flex;
  align-items: center;
  gap: 0.5em;
  color: var(--textSecondary, #888);
  font-size: 0.9em;
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 0.4em;
  color: #ef4444;
  font-size: 0.85em;
}
</style>

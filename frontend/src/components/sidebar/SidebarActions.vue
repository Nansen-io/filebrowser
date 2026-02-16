<template>
  <transition name="expand" @before-enter="beforeEnter" @enter="enter" @leave="leave">
    <div v-if="showActions" class="sidebar-actions card">
      <div class="inner-card">
        <!-- New Folder -->
        <button
          v-if="canCreate"
          @click="handleCreateFolder"
          class="action button action-button"
          :aria-label="$t('files.newFolder')"
        >
          <i class="material-icons action-icon">create_new_folder</i>
          <span>{{ $t('files.newFolder') }}</span>
        </button>

        <!-- New File -->
        <button
          v-if="canCreate"
          @click="handleCreateFile"
          class="action button action-button"
          :aria-label="$t('files.newFile')"
        >
          <i class="material-icons action-icon">note_add</i>
          <span>{{ $t('files.newFile') }}</span>
        </button>

        <!-- Upload -->
        <button
          v-if="canCreate"
          @click="handleUpload"
          class="action button action-button"
          :aria-label="$t('general.upload')"
        >
          <i class="material-icons action-icon">file_upload</i>
          <span>{{ $t('general.upload') }}</span>
        </button>

        <!-- Share -->
        <button
          v-if="canShare"
          @click="handleShare"
          class="action button action-button"
          :aria-label="$t('general.share')"
        >
          <i class="material-icons action-icon">share</i>
          <span>{{ $t('general.share') }}</span>
        </button>
      </div>
    </div>
  </transition>
</template>

<script>
import { state, getters } from "@/store";
import { useFileActions } from "@/composables/useFileActions";

export default {
  name: "SidebarActions",
  setup() {
    const {
      createNewFolder,
      createNewFile,
      uploadFiles,
      shareCurrentFolder
    } = useFileActions();

    return {
      createNewFolder,
      createNewFile,
      uploadFiles,
      shareCurrentFolder
    };
  },
  computed: {
    permissions() {
      return getters.permissions();
    },
    isShare() {
      return getters.isShare();
    },
    isSearchActive() {
      return state.isSearchActive;
    },
    isListingView() {
      return getters.currentView() === "listingView";
    },
    canCreate() {
      return this.permissions.create && !this.isShare && !this.isSearchActive;
    },
    canShare() {
      return this.permissions.share && !this.isShare;
    },
    showActions() {
      return this.isListingView && (this.canCreate || this.canShare);
    }
  },
  methods: {
    handleCreateFolder() {
      this.createNewFolder();
    },
    handleCreateFile() {
      this.createNewFile();
    },
    handleUpload() {
      this.uploadFiles();
    },
    handleShare() {
      this.shareCurrentFolder();
    },
    beforeEnter(el) {
      el.style.height = '0';
      el.style.opacity = '0';
    },
    enter(el, done) {
      el.style.transition = '';
      el.style.height = '0';
      el.style.opacity = '0';
      void el.offsetHeight;
      el.style.transition = 'height 0.3s, opacity 0.3s';
      el.style.height = el.scrollHeight + 'px';
      el.style.opacity = '1';
      setTimeout(() => {
        el.style.height = 'auto';
        done();
      }, 300);
    },
    leave(el, done) {
      el.style.transition = 'height 0.3s, opacity 0.3s';
      el.style.height = el.scrollHeight + 'px';
      void el.offsetHeight;
      el.style.height = '0';
      el.style.opacity = '0';
      setTimeout(done, 300);
    }
  }
};
</script>

<style scoped>
.sidebar-actions {
  padding: 1em;
  margin-top: 0.5em;
}

.inner-card {
  display: flex;
  flex-direction: column;
  gap: 0.5em;
  width: 100%;
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  width: 100%;
  padding: 0.75em 1em;
  border-radius: 0.5em;
  background-color: var(--background);
  color: var(--textPrimary);
  border: none;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.1s;
  gap: 0.75em;
}

.action-button:hover {
  background-color: var(--surfaceSecondary);
  transform: translateY(-2px);
}

.action-button:active {
  transform: translateY(0);
}

.action-button span {
  color: var(--textPrimary);
}

.action-icon {
  color: var(--textPrimary);
  font-size: 1.25em;
}

/* Animation styles */
.expand-enter-active,
.expand-leave-active {
  transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.expand-enter,
.expand-leave-to {
  height: 0 !important;
  opacity: 0;
}
</style>

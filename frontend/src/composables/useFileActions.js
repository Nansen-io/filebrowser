import { mutations, state } from "@/store";

/**
 * Composable for file actions that can be reused across components
 * Provides consistent action handlers for file operations
 */
export function useFileActions() {
  const createNewFolder = () => {
    mutations.showHover("newDir");
  };

  const createNewFile = () => {
    mutations.showHover("newFile");
  };

  const uploadFiles = () => {
    mutations.showHover("upload");
  };

  const shareCurrentFolder = () => {
    mutations.showHover({
      name: "share",
      props: { item: state.req }
    });
  };

  return {
    createNewFolder,
    createNewFile,
    uploadFiles,
    shareCurrentFolder
  };
}

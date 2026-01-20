# Frontend

Display contents of right click context menu in side bar.

I want to add this section which is the right click context menu:

```html
<div data-v-d98731eb="" id="context-menu" style="top: 425px; left: 573px; height: 250px; opacity: 1; visibility: visible; transition: height 0.3s, opacity 0.3s;" class="button no-select fb-shadow dark-mode"><div data-v-d98731eb="" class="context-menu-header"><div data-v-d98731eb="" class="action button clickable"><!----><i data-v-d98731eb="" class="material-icons">arrow_back</i></div><!----></div><hr data-v-d98731eb="" class="divider"><button data-v-d98731eb="" aria-label="New folder" title="New folder" class="action no-select"><i class="material-icons">create_new_folder</i><span>New folder</span><!----></button><button data-v-d98731eb="" aria-label="New file" title="New file" class="action no-select"><i class="material-icons">note_add</i><span>New file</span><!----></button><button data-v-d98731eb="" aria-label="Upload" title="Upload" class="action no-select"><i class="material-icons">file_upload</i><span>Upload</span><!----></button><!----><!----><button data-v-d98731eb="" aria-label="Share" title="Share" class="action no-select"><i class="material-icons">share</i><span>Share</span><!----></button><!----><!----><!----><!----><!----><!----><!----><!----></div>
```

also see Capture.PNG for refence if needed on right click context menu.

I want to add the exact same options to the righthand bar on the main display blew the "Links" Menu which is

```html
<div data-v-408f0343="" class="sidebar-links card"><div data-v-408f0343="" class="sidebar-links-header"><i data-v-408f0343="" class="material-icons action">home</i><span data-v-408f0343="">Links</span><i data-v-408f0343="" class="material-icons action">edit</i></div><div data-v-408f0343="" class="inner-card"><a data-v-408f0343="" href="/files/srv/" class="action button source-button sidebar-link-button active" aria-label="srv"><div data-v-408f0343="" class="source-container"><svg data-v-408f0343="" class="realtime-pulse ready"><circle data-v-408f0343="" class="center" cx="50%" cy="50%" r="7px"></circle><circle data-v-408f0343="" class="pulse" cx="50%" cy="50%" r="10px"></circle></svg><span data-v-408f0343="">srv</span><i data-v-408f0343="" class="no-select material-symbols-outlined tooltip-info-icon"> info </i></div><div data-v-408f0343="" class="usage-info"><div data-v-408f0343="" class="vue-simple-progress" style="background: rgb(238, 238, 238); position: relative; min-height: 16px; border-radius: 8px;"><!----><div class="vue-simple-progress-bar" style="width: 0%; height: 16px; background: var(--primaryColor); transition: 0.5s; border-radius: 8px;"></div><div class="vue-simple-progress-text" style="color: rgb(0, 0, 0); font-size: 13px; text-align: center; position: absolute; left: 0px; right: 0px; top: 50%; transform: translateY(-50%); width: 100%; padding: 0px 0.5em; box-sizing: border-box;">143.5 MB / 953.2 GB (0%) <!----></div></div></div></a></div></div>
```

It will display the same options as the right click context menu which are:

- New Folder
- New file
- Upload
- Share

# User

Display Username not UUID
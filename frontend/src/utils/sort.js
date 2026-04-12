function getExtension(name) {
    if (!name || !name.includes('.')) return '';
    const ext = name.split('.').pop();
    return ext === name ? '' : ext.toLowerCase();
}

export function sortedItems(items = [], sortby="name", asc=true) {
    return items.sort((a, b) => {
        let valueA = a[sortby];
        let valueB = b[sortby];

        // Special handling for protected status (true sorts first when ascending)
        if (sortby === "protected") {
            const aVal = a.protected ? 1 : 0;
            const bVal = b.protected ? 1 : 0;
            return asc ? bVal - aVal : aVal - bVal;
        }

        // Special handling for duration which is stored in metadata
        if (sortby === "duration") {
            valueA = a.metadata?.duration ?? 0;
            valueB = b.metadata?.duration ?? 0;
        }

        // Special handling for extension/type sorting
        if (sortby === "extension") {
            valueA = a.type === "directory" ? "\x00" : getExtension(a.name);
            valueB = b.type === "directory" ? "\x00" : getExtension(b.name);
            const comparison = valueA.localeCompare(valueB, undefined, { sensitivity: 'base' });
            return asc ? comparison : -comparison;
        }

        if (sortby === "name") {
            // Use localeCompare with numeric option for natural sorting
            const comparison = valueA.localeCompare(valueB, undefined, { numeric: true, sensitivity: 'base' });
            return asc ? comparison : -comparison;
        }

        // Default sorting for other fields (including created, modified — ISO strings compare correctly)
        if (asc) {
            return valueA > valueB ? 1 : -1;
        } else {
            return valueA < valueB ? 1 : -1;
        }
    });
}

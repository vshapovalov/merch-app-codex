<template>
  <section class="page">
    <header class="flex align-items-center justify-content-between mb-4">
      <h2 class="m-0">Категории</h2>
      <Button label="Добавить корневую категорию" icon="pi pi-plus" @click="openCreateRoot" />
    </header>

    <Tree
      v-if="categoriesTree.length"
      :value="categoriesTree"
      :loading="loading"
      selectionMode="single"
      class="categories-tree"
    >
      <template #default="{ node }">
        <div class="tree-node flex align-items-center justify-content-between w-full">
          <span class="tree-node-label">{{ node.label }}</span>
          <div class="tree-node-actions flex gap-2">
            <Button
              icon="pi pi-plus"
              severity="success"
              text
              rounded
              @click.stop="openCreateChild(node.data.id)"
            />
            <Button icon="pi pi-pencil" severity="info" text rounded @click.stop="openEdit(node.data)" />
            <Button icon="pi pi-trash" severity="danger" text rounded @click.stop="deleteItem(node.data)" />
          </div>
        </div>
      </template>
    </Tree>
    <div v-else-if="!loading" class="empty-placeholder">Категории пока не добавлены</div>

    <Dialog v-model:visible="dialogVisible" modal :header="dialogTitle" class="w-full md:w-6">
      <form v-if="currentItem" class="flex flex-column gap-3" @submit.prevent="saveItem">
        <span class="p-float-label">
          <InputText id="name" v-model="currentItem.name" required />
          <label for="name">Название</label>
        </span>
        <div>
          <label class="block mb-2">Родительская категория</label>
          <Dropdown
            v-model="currentItem.parent_id"
            :options="parentOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="Не выбрана"
            showClear
            class="w-full"
          />
        </div>
        <div class="flex justify-content-end gap-2">
          <Button label="Отмена" type="button" severity="secondary" text @click="dialogVisible = false" />
          <Button label="Сохранить" type="submit" :loading="saving" />
        </div>
      </form>
    </Dialog>
  </section>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useCrud } from '../composables/useCrud';

const {
  items,
  loading,
  saving,
  dialogVisible,
  currentItem,
  loadItems,
  openCreate,
  openEdit,
  saveItem,
  deleteItem,
} = useCrud('/categories', () => ({ name: '', parent_id: null }));

const dialogTitle = computed(() =>
  currentItem?.value?.id ? 'Редактирование категории' : 'Новая категория'
);

const parentOptions = computed(() => {
  const currentId = currentItem?.value?.id;
  return items.value.filter((category) => category.id !== currentId);
});

const buildTree = (categories, parentId = null) => {
  return categories
    .filter((category) => category.parent_id === parentId)
    .sort((a, b) => a.name.localeCompare(b.name))
    .map((category) => ({
      key: String(category.id),
      label: category.name,
      data: category,
      expanded: true,
      children: buildTree(categories, category.id),
    }));
};

const categoriesTree = computed(() => buildTree(items.value));

const openCreateRoot = () => openCreate({ parent_id: null });

const openCreateChild = (parentId) => openCreate({ parent_id: parentId });

onMounted(loadItems);
</script>

<style scoped>
.page {
  background: var(--surface-card);
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: var(--card-shadow, 0 1px 4px rgba(0, 0, 0, 0.1));
}

.categories-tree {
  border: 1px solid var(--surface-border);
  border-radius: 0.75rem;
  padding: 1rem;
}

.tree-node-label {
  font-weight: 500;
}

.tree-node-actions .p-button {
  width: 2rem;
  height: 2rem;
}

.empty-placeholder {
  text-align: center;
  color: var(--text-color-secondary);
  padding: 2rem 0;
}
</style>

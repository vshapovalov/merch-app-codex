import { ref } from 'vue';
import { useToast } from 'primevue/usetoast';
import api from '../services/api';

export function useCrud(endpoint, createDefault, options = {}) {
  const { preparePayload } = options;
  const items = ref([]);
  const loading = ref(false);
  const saving = ref(false);
  const dialogVisible = ref(false);
  const currentItem = ref(null);
  const toast = useToast();

  const loadItems = async () => {
    loading.value = true;
    try {
      const { data } = await api.get(endpoint);
      items.value = data;
    } catch (error) {
      console.error(error);
      toast.add({ severity: 'error', summary: 'Ошибка', detail: 'Не удалось загрузить данные', life: 3000 });
    } finally {
      loading.value = false;
    }
  };

  const openCreate = () => {
    currentItem.value = createDefault();
    dialogVisible.value = true;
  };

  const openEdit = (item) => {
    currentItem.value = { ...item };
    dialogVisible.value = true;
  };

  const buildPayload = () => {
    if (!currentItem.value) {
      return {};
    }
    return preparePayload ? preparePayload({ ...currentItem.value }) : currentItem.value;
  };

  const saveItem = async () => {
    saving.value = true;
    try {
      const payload = buildPayload();
      if (currentItem.value.id) {
        await api.put(`${endpoint}/${currentItem.value.id}`, payload);
        toast.add({ severity: 'success', summary: 'Сохранено', detail: 'Запись обновлена', life: 2000 });
      } else {
        const { data } = await api.post(endpoint, payload);
        currentItem.value = data;
        toast.add({ severity: 'success', summary: 'Создано', detail: 'Запись добавлена', life: 2000 });
      }
      await loadItems();
      dialogVisible.value = false;
    } catch (error) {
      console.error(error);
      toast.add({ severity: 'error', summary: 'Ошибка', detail: 'Не удалось сохранить запись', life: 3000 });
    } finally {
      saving.value = false;
    }
  };

  const deleteItem = async (item) => {
    try {
      await api.delete(`${endpoint}/${item.id}`);
      toast.add({ severity: 'success', summary: 'Удалено', detail: 'Запись удалена', life: 2000 });
      await loadItems();
    } catch (error) {
      console.error(error);
      toast.add({ severity: 'error', summary: 'Ошибка', detail: 'Не удалось удалить запись', life: 3000 });
    }
  };

  return {
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
  };
}

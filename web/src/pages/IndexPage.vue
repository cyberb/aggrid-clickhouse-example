<template>
  <q-page padding style="height:100px" >
    <ag-grid-vue
      class="ag-theme-alpine-dark fit"
      :columnDefs="columnDefs"
      @grid-ready="onGridReady"
      :defaultColDef="defaultColDef"
      :rowModelType="rowModelType"
      :enableCellChangeFlash="true"
      :suppressAggFuncInHeader="true"
      :animateRows="true"
      debounceVerticalScrollbar=true
      :debug="true"
      :maxBlocksInCache="0"
      :getRowId="getRowId"
      @row-group-opened="onColumnRowGroupChanged"
    ></ag-grid-vue>
  </q-page>
</template>

<script lang="ts">
import "ag-grid-community/styles/ag-grid.css";
import "ag-grid-community/styles/ag-theme-alpine.css";
import 'ag-grid-enterprise';
import {AgGridVue} from "ag-grid-vue3";
import {defineComponent, onBeforeMount, ref} from "vue";
import axios from "axios";


export default defineComponent({
    name: 'IndexPage',
    components: {
      AgGridVue,
    },
    setup() {
      const columnDefs = ref([
        {
          field: "id", minWidth: 300, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true
        },

        {
          field: 'year', minWidth: 100, sortable: true,
          filter: "agNumberColumnFilter", floatingFilter: true, rowGroup: "true"
        },
        {
          field: "amount", minWidth: 200, sortable: true,
          filter: "agNumberColumnFilter", floatingFilter: true, aggFunc: "sum"
        },
        {
          field: "field1", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },
        {
          field: "field2", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },
        {
          field: "field3", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },
        {
          field: "field4", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },
        {
          field: "field5", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },
        {
          field: "field6", minWidth: 100, sortable: true,
          filter: "agTextColumnFilter", floatingFilter: true,
        },

      ])
      const gridApi = ref();
      const gridColumnApi = ref();
      const defaultColDef = ref({
        flex: 1,
        minWidth: 100,
      })

      const rowModelType = ref();

      onBeforeMount(() => {
        rowModelType.value = "serverSide"
      })

      const createServerSideDatasource = () => {
        return {
          getRows: (params: any) => {
            console.debug("get rows ", params)
            axios.post('/api/clickhouse', params.request)
              .then((response: any) => {
                params.success(response.data);
              })
              .catch((err: any) => {
                console.error(Error(err));
                params.fail();
              })
          }
        }
      }

      const onGridReady = (params: any) => {
        console.debug("onGridReady")

        gridApi.value = params.api;
        gridColumnApi.value = params.columnApi;

        params.api.blockLoadDebounceMillis = 1000
        params.api.maxBlocksInCache = 0;

        params.api.setServerSideDatasource(createServerSideDatasource());
      }
      const getRowId = (params: any) => {
        if (params.parentKeys) {
          return params.parentKeys[0] + '-' + params.data.id
        }

        return params.data.year;
      }
      const groups: Array<string> = []

      const refreshCache = () => {
        gridApi.value.refreshServerSide({purge: false});
        groups.forEach((group) => {
          console.debug("refresh group: " + group)

          gridApi.value.refreshServerSide({route: [group], purge: false});

        })
      }

      const onRowGroupOpened = (event: any) => {
        if (event.expanded) {
          groups.push(event.node.key)
        } else {
          delete groups[event.node.key]
        }
        console.debug(event)
      }

      setInterval(() => {
        if (gridApi.value.destroyCalled) {
          return;
        }
        refreshCache()
      }, 1000)

      return {
        columnDefs,
        gridApi,
        gridColumnApi,
        defaultColDef,
        rowModelType,
        onGridReady,
        getRowId,
        onColumnRowGroupChanged: onRowGroupOpened,
      }
    },
  }
)
</script>

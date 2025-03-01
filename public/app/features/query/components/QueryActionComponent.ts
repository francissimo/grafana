import { DataQuery, TimeRange } from '@grafana/data';

interface ActionComponentProps {
  query?: DataQuery;
  queries?: Array<Partial<DataQuery>>;
  onAddQuery?: (q: DataQuery) => void;
  timeRange?: TimeRange;
}

type QueryActionComponent = React.ComponentType<ActionComponentProps>;

class QueryActionComponents {
  extraRenderActions: QueryActionComponent[] = [];

  addExtraRenderAction(extra: QueryActionComponent) {
    this.extraRenderActions = this.extraRenderActions.concat(extra);
  }

  getAllExtraRenderAction(): QueryActionComponent[] {
    return this.extraRenderActions;
  }
}

/**
 * @internal and experimental
 */
export const GroupActionComponents = new QueryActionComponents();

/**
 * @internal and experimental
 */
export const RowActionComponents = new QueryActionComponents();

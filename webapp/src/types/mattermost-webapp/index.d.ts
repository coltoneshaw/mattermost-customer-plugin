export type PluginComponent = {
    id: string;
    pluginId: string;
    title?: string;

    /** @default null - which means 'channels'*/
    supportedProductIds?: ProductScope;
    component?: React.ComponentType;
    subMenu?: Menu[];
    text?: string;
    dropdownText?: string;
    tooltipText?: string;
    button?: React.ReactElement;
    dropdownButton?: React.ReactElement;
    icon?: React.ReactElement;
    iconUrl?: string;
    mobileIcon?: React.ReactElement;
    filter?: (id: string) => boolean;
    action?: (...args: unknown) => void; // TODO Add more concrete types?
    shouldRender?: (state: GlobalState) => boolean;
    hook?: (post: Post, message?: string) => string;
};

type ReactResolvable = React.ReactNode | React.ElementType;

export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType);

    /**
    * Register a Right-Hand Sidebar component by providing a title for the right hand component.
    * Accepts the following:
    * - component - A react component to display in the Right-Hand Sidebar.
    * - title - A string or JSX element to display as a title for the RHS.
    * Returns:
    * - id: a unique identifier
    * - showRHSPlugin: the action to dispatch that will open the RHS.
    * - hideRHSPlugin: the action to dispatch that will close the RHS
    * - toggleRHSPlugin: the action to dispatch that will toggle the RHS
    */
    registerRightHandSidebarComponent(component: PluginComponent['component'], title: string): {
        id: string;
        showRHSPlugin: Action<Record<string, unknown>>;
        hideRHSPlugin: Action<Record<string, unknown>>;
        toggleRHSPlugin: Action<Record<string, unknown>>;
    };

    /**
     * INTERNAL: Subject to change without notice.
     * Add an item to the App Bar.
     * @param {string} iconUrl resolvable URL to use as the button's icon.
     * @param {PluginComponent['action'] | undefined} action called when the button is clicked, passed the channel and channel member as arguments.
     * @param {React.ReactNode} tooltipText string or React element shown for tooltip appear on hover.
     * @param {null | string | Array<null | string>} supportedProductIds specifies one or multiple product identifier(s),
     * identifiers can either be the "real" product uuid, or a product's more commonly accessible plugin id, or '*' to match everything.
     * @param {PluginComponent['component'] | undefined} rhsComponent an optional corresponding RHS component. If provided, its toggler is automatically wired to the action.
     * @param {ReactResolvable | undefined} rhsTitle the corresponding RHS component's title.
     * @returns {string} unique identifier
     */
    registerAppBarComponent({
        iconUrl,
        action,
        tooltipText,
        supportedProductIds = null,
        rhsComponent,
        rhsTitle,
    }: {
        iconUrl: string;
        tooltipText: ReactResolvable;
        supportedProductIds: ProductScope;
    } & ({
        action: PluginComponent['action'];
        rhsComponent?: never;
        rhsTitle?: never;
    } | {
        action?: never;
        rhsComponent: PluginComponent;
        rhsTitle: ReactResolvable;
    })) : string | {
        id: string;
        component: ReturnType<PluginRegistry['registerRightHandSidebarComponent']>
    }

    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}

import * as cog from '../cog';
import * as builder_delegation from '../builder_delegation';

export class DashboardBuilder implements cog.Builder<builder_delegation.Dashboard> {
    private readonly internal: builder_delegation.Dashboard;

    constructor() {
        this.internal = builder_delegation.defaultDashboard();
    }

    build(): builder_delegation.Dashboard {
        return this.internal;
    }

    id(id: number): this {
        this.internal.id = id;
        return this;
    }

    title(title: string): this {
        this.internal.title = title;
        return this;
    }

    links(links: cog.Builder<builder_delegation.DashboardLink>[]): this {
        const linksResources = links.map(builder => builder.build());
        this.internal.links = linksResources;
        return this;
    }

    singleLink(singleLink: cog.Builder<builder_delegation.DashboardLink>): this {
        const singleLinkResource = singleLink.build();
        this.internal.singleLink = singleLinkResource;
        return this;
    }
}
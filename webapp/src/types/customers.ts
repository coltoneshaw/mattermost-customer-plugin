import {AdminConfig} from '@mattermost/types/config';

// eslint-disable-next-line no-shadow
export const enum LicenseType {
    Cloud = 'cloud',
    Enterprise = 'enterprise',
    Professional = 'professional',
    Free = 'free',
    Trial = 'trial',
    Nonprofit = 'nonprofit',
    Other = 'other'
}

export type Customer = {
    id: string;
    name: string;
    lastUpdated: number;
    salesforceId: string;
    zendeskId: string;
    customerSuccessManager: string;
    accountExecutive: string;
    technicalAccountManager: string;
    productManager: string;
    licensedTo: string;
    siteURL: string;
    licenseType: LicenseType | '';
    customerChannel: string;
    GDriveLink: string;
    airGapped: boolean;
    airGappedReason: string;
    region: string;
    status: string;
    companyType: string;
    codeWord: string;
}

export type CustomerPacketValues = {
    licensedTo: string;
    version: string;
    serverOS: string;
    serverArch: string;
    databaseType: string;
    databaseVersion: string;
    databaseSchemaVersion: string;
    fileDriver: string;
    activeUsers: number;
    dailyActiveUsers: number;
    monthlyActiveUsers: number;
    inactiveUserCount: number;
    licenseSupportedUsers: number;
    totalPosts: number;
    totalChannels: number;
    totalTeams: number;
    elasticServerVersion: string;
    metrics: boolean;
    metricService: string;
    hostingType: string;
    deploymentType: string;
    mobileApp: string;
    productsInUse: string;
    samlProvider: string;
    ldapProvider: string;
}

export type CustomerPluginValues = {
    homePageURL: string;
    pluginID: string;
    version: string;
    isActive: boolean,
    name: string
}

export type CustomerConfigValues = AdminConfig;

export type FullCustomerInfo = Customer & {
    packet: CustomerPacketValues;
    config: AdminConfig;
    plugins: CustomerPluginValues[]
}

export type GetCustomerResult = {
    totalCount: number;
    pageCount: number;
    hasMore: boolean;
    customers: FullCustomerInfo[]
}

// eslint-disable-next-line no-shadow
export enum CustomerSortOptions {
    SortByName = 'name',
    SortByCSM = 'csm',
    SortByAE = 'ae',
    SortByTAM = 'tam',
    SortByType = 'type',
    SortBySiteURL = 'site_url',
    SortByLicensedTo = 'licensed_to',
    SortByLastUpdated = 'last_updated',
    Default = ''
}

// eslint-disable-next-line no-shadow
export enum SortDirection {
    DirectionAsc = 'asc',
    DirectionDesc = 'desc'
}

export type CustomerFilterOptions = {
    sort: CustomerSortOptions;
    order: SortDirection;
    page: string;
    perPage: string;
    searchTerm: string;
}

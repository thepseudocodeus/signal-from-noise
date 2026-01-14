export namespace main {
	
	export class CreateZipRequest {
	    production_request_id: string;
	    file_ids: number[];
	
	    static createFrom(source: any = {}) {
	        return new CreateZipRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.production_request_id = source["production_request_id"];
	        this.file_ids = source["file_ids"];
	    }
	}
	export class CreateZipResponse {
	    zip_path: string;
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateZipResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.zip_path = source["zip_path"];
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class FileInfo {
	    id: number;
	    path: string;
	    directory: string;
	    category: string;
	    date: string;
	    size: number;
	    privileged: boolean;
	    duplicate_hash: string;
	    file_name: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.directory = source["directory"];
	        this.category = source["category"];
	        this.date = source["date"];
	        this.size = source["size"];
	        this.privileged = source["privileged"];
	        this.duplicate_hash = source["duplicate_hash"];
	        this.file_name = source["file_name"];
	    }
	}
	export class FileSearchRequest {
	    production_request_id: string;
	    date_start: string;
	    date_end: string;
	    categories: string[];
	    exclude_privileged: boolean;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new FileSearchRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.production_request_id = source["production_request_id"];
	        this.date_start = source["date_start"];
	        this.date_end = source["date_end"];
	        this.categories = source["categories"];
	        this.exclude_privileged = source["exclude_privileged"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class FileSearchResult {
	    files: FileInfo[];
	    total_count: number;
	    page: number;
	    page_size: number;
	    total_pages: number;
	
	    static createFrom(source: any = {}) {
	        return new FileSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], FileInfo);
	        this.total_count = source["total_count"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	        this.total_pages = source["total_pages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


export namespace main {
	
	export class ProductionRequest {
	    id: number;
	    title?: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	    }
	}

}


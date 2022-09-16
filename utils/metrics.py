import requests
# def factory_metric():
#         return {'key': None, 'label':[], 'value': 0}
# res = defaultdict(factoryheader)
class Metrics(object):
    """
    Peer Metrics 解释器
    """

    def __init__(self, url, filters:list):
        """
        param: filter: metrics参数过滤器
        """
        self._filter = filters
        self._url = url
    
    @property
    def schema(self):
        return {'key': None, 'label':{}, 'value': 0}
        

    def interprete(self) -> dict:
        """
        解析
        """
        metrics = []
        
        response = requests.request("GET", self._url)
        data = response.text

        lines = data.split('\n')
        for line in lines:
            if line.startswith('# ') or len(line) == 0:
                continue
 
            metric_item = self.schema
            line_parts = line.split(' ')
            metric_item['value'] = line_parts[1]
            
            if line_parts[0].endswith('}'):
                metric_item['key'] = line_parts[0].split('{')[0]
                labels = line_parts[0].split('{')[1][:-1].split(',')
                for la in labels:
                    metric_item['label'][la.split('=')[0]] = la.split('=')[1][1:-1]
            else:
                metric_item['key'] = line_parts[0]

            metrics.append(metric_item)

        # filter
        if self._filter:
            result = []
            for item in metrics:
                if item['key'] in self._filter:
                    result.append(item)
            
            return result

        return metrics




    
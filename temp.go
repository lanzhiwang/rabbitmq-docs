func (r *RabbitmqClusterReconciler) GetRabbitmqLogs(ctx context.Context, namespace string, podName string, opts *corev1.PodLogOptions) (([]string, error)) {
	RabbitmqLogs, err := corev1client.NewForConfig(r.ClusterConfig).Pods(namespace).GetLogs(podName, opts).Stream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get pod logs stream")
	}
	defer RabbitmqLogs.Close()

	logArr := make([]string, 0)
	sc := bufio.NewScanner(RabbitmqLogs)
	for sc.Scan() {
		logArr = append(logArr, sc.Text())
	}
	return logArr, errors.Wrap(sc.Err(), "reading logs stream")
}

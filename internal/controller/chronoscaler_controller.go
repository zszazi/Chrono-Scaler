package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/zszazi/Chrono-Scaler/api/v1alpha1"
)

var logger = log.Log.WithName("controller_chrono-scaler")

// ChronoScalerReconciler reconciles a ChronoScaler object
type ChronoScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.zszazi.github.io,resources=chronoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.zszazi.github.io,resources=chronoscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.zszazi.github.io,resources=chronoscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ChronoScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	log.Info("Reconciler Loop Running...")
	scaler := &apiv1alpha1.ChronoScaler{}

	if err := r.Get(ctx, req.NamespacedName, scaler); err != nil {
		return ctrl.Result{}, nil
	}
	startTime := scaler.Spec.Start
	endTime := scaler.Spec.End
	replicas := scaler.Spec.Replicas
	defaultReplicas := scaler.Spec.DefaultReplicas

	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	currentTime := time.Since(midnight)

	start, err := parseDurationFromTimeString(startTime)
	if err != nil {
		log.Info("Error parsing start time:")
		return ctrl.Result{}, nil
	}

	end, err := parseDurationFromTimeString(endTime)
	if err != nil {
		log.Info("Error parsing end time:")
		return ctrl.Result{}, nil
	}

	var targetReplicas int32
	if start <= currentTime && currentTime <= end {
		targetReplicas = replicas
	} else if currentTime > end {
		targetReplicas = defaultReplicas
	}

	for _, deploy := range scaler.Spec.Deployments {
		deployment := &v1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: deploy.Namespace,
			Name:      deploy.Name,
		}, deployment)

		if err != nil {
			return ctrl.Result{}, err
		}

		if *deployment.Spec.Replicas != targetReplicas {
			targetReplicasStr := strconv.Itoa(int(targetReplicas))
			currentReplicasStr := strconv.Itoa(int(*deployment.Spec.Replicas))

			log.Info("Change in the replicas detected",
				"current_time", currentTime,
				"scaling_to", targetReplicasStr,
				"current_replicas", currentReplicasStr)

			deployment.Spec.Replicas = &targetReplicas
			err := r.Update(ctx, deployment)
			if err != nil {
				scaler.Status.Status = apiv1alpha1.FAILED
				return ctrl.Result{}, err
			}
			scaler.Status.Status = apiv1alpha1.SUCCESS
			err = r.Status().Update(ctx, scaler)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func parseDurationFromTimeString(timeStr string) (time.Duration, error) {
	parts := strings.Split(timeStr, "h")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour value: %v", err)
	}

	minutesStr := strings.TrimSuffix(parts[1], "m")
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return 0, fmt.Errorf("invalid minute value: %v", err)
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChronoScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.ChronoScaler{}).
		Complete(r)
}
